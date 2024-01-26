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
}

func NewServer(c ...ServerConfigFunc) *Server {
	s := &Server{}
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
	http.HandleFunc("/terraform-state", handleTerraformState)
}

func handleTerraformState(w http.ResponseWriter, r *http.Request) {
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
