package main

import (
	s "abcode.com/vm-code-service/server"
)

func main() {
	server := s.NewServer()
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
