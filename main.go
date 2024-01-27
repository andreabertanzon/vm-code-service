package main

import (
	s "abcode.com/vm-code-service/server"
	"abcode.com/vm-code-service/services"
)

func main() {
	minioService, err := services.NewMinioService()
	if err != nil {
		panic(err)
	}
	server := s.NewServer(minioService)

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
