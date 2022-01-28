package main

import (
	"fmt"
	"net/http"
)

func StartServer() {
	apiServerMux := http.NewServeMux()
	apiServer := http.Server {
	Addr: fmt.Sprintf(":%v", 8080),
	Handler: apiServerMux,
	}
	apiServerMux.HandleFunc("/api", handler)
	apiServerMux.HandleFunc("/api/returnMenu", handler)
	apiServerMux.HandleFunc("/api/returnAccounts", handler)
	apiServerMux.HandleFunc("/api/prepareCourse", handler)
	apiServerMux.HandleFunc("/api/serveCourse", handler)
	PrintPositive("Kitchen is ready.")
	CheckForError(apiServer.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w)
}