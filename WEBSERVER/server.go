package main

import (
	"fmt"
	"./services"
	"golang.org/x/net/http2"
	"log"
	"net/http"
)

func main() {
	// redirect every http request to https
	go http.ListenAndServe(services.Config.Address, http.HandlerFunc(services.Redirect))
	
	mux := http.NewServeMux()

	// serving style and js files(only if no Bootstrap CDM):
	statics := http.FileServer(http.Dir(services.Config.Static))
	mux.Handle("/cssjs/", http.StripPrefix("/cssjs/", statics))

	// defined in server_routes.go
	mux.HandleFunc("/", services.Index)
	mux.HandleFunc("/err", services.Err)
	mux.HandleFunc("/about", services.About)


	// javascript calls:
	mux.HandleFunc("/bruse", services.Bruse)

	// // defined in route_thread.go
	// mux.HandleFunc("/thread/new", go_services.NewThread)
	// mux.HandleFunc("/thread/create", go_services.CreateThread)
	// mux.HandleFunc("/thread/post", go_services.PostThread)
	// mux.HandleFunc("/thread/read", go_services.ReadThread)

	server := http.Server{
		Addr:    services.Config.AddressSSL,
		Handler: mux,
	}
	fmt.Println(services.Config.Address + " and " + services.Config.AddressSSL)
	http2.ConfigureServer(&server, &http2.Server{})
	log.Fatal(server.ListenAndServeTLS("/etc/letsencrypt/live/datapenetration.de/fullchain.pem", "/etc/letsencrypt/live/datapenetration.de/privkey.pem"))
}
