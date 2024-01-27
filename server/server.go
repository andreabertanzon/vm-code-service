package server

import (
	"fmt"
	"net/http"

	"abcode.com/vm-code-service/services"
)

type ServerConfigFunc func(c *ServerConfig)
type ServerConfig struct {
	port        string
	environment string
	tls         bool
	configFile  string
}

func defaultServerConfig() ServerConfig {
	return ServerConfig{
		port:        "8080",
		environment: "development",
		tls:         false,
		configFile:  "config.yaml",
	}
}

type Server struct {
	config ServerConfig
	services.S3Service
}

func NewServer(tfHandler services.S3Service, c ...ServerConfigFunc) *Server {
	s := &Server{
		S3Service: tfHandler,
	}
	conf := defaultServerConfig()
	for _, configFn := range c {
		configFn(&conf)
	}
	return s
}

func (s *Server) Start() error {
	s.registerRoutes()
	err := http.ListenAndServe(s.config.port, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) registerRoutes() {
	http.HandleFunc("/terraform-state", s.handleTerraformState)
	http.HandleFunc("/template-content", s.handleDowloadZipFolder)
}

func (s *Server) handleTerraformState(w http.ResponseWriter, r *http.Request) {
	// handle the files by query parameter like file=

	minioService, err := services.NewMinioService()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error creating minio service")
		return
	}
	outobj, err := minioService.GetTerraformState()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error getting terraform state")
		return
	}
	w.WriteHeader(200)
	w.Header().Add("Content-Disposition", "attachment; filename="+"terraform.tfstate")
	w.Write(outobj)
}

func (s *Server) handleDowloadZipFolder(w http.ResponseWriter, r *http.Request) {
	// only handle GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		fmt.Fprint(w, "Method not allowed")
		return
	}

	folderQueryParameter := r.URL.Query().Get("folder")

	content, err := s.S3Service.DowloadBucketFolderToZip("vm-templates", folderQueryParameter)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error getting terraform state")
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Disposition", "attachment; filename="+folderQueryParameter+".zip")
	w.Write(content)
}
