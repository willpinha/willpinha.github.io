package main

import (
	"log"
	"net/http"
	"strconv"
)

func runServer(dir string, port int) error {
	addr := ":" + strconv.Itoa(port)
	log.Printf("serving %s on http://localhost%s", dir, addr)
	return http.ListenAndServe(addr, http.FileServer(http.Dir(dir)))
}
