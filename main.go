package main

import (
	"log"
	"net/http"
	swagger "oldcow/go"
	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//
)

const (
	CertPath string = "candy.tld/cert.pem"
	KeyPath  string = "candy.tld/key.pem"
)

func main() {
	log.Printf("Server started")
	router := swagger.NewRouter()

	log.Fatal(http.ListenAndServeTLS(":3333", CertPath, KeyPath, router))
}
