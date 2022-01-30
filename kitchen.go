package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func StartServer() {
	apiServerMux := http.NewServeMux()
	apiServer := http.Server {
	Addr: fmt.Sprintf(":%v", 8080),
	Handler: apiServerMux,
	}
	apiServerMux.HandleFunc("/api", handler)
	apiServerMux.HandleFunc("/api/returnMenu", returnMenu)
	apiServerMux.HandleFunc("/api/addCourse", addCourse)
	apiServerMux.HandleFunc("/api/removeCourse/{id}", removeCourse)
	PrintPositive("Kitchen is ready.")
	CheckForError(apiServer.ListenAndServe())
}

func handler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println(w)
}

func returnMenu(w http.ResponseWriter, _ *http.Request) {
	setWriter(w)
	dishes := GetDishes()
	CheckForError(json.NewEncoder(w).Encode(dishes))
}

func addCourse(w http.ResponseWriter, r *http.Request) {
	setWriter(w)
	var courseRequest Dish
	err := json.NewDecoder(r.Body).Decode(&courseRequest)
	CheckForError(err)
	AddDish(courseRequest)
	w.WriteHeader(http.StatusOK)
}

func removeCourse(w http.ResponseWriter, r *http.Request) {
	setWriter(w)
	dishId := r.URL.Path[strings.LastIndex(r.URL.Path, "/") + 1:]
	RemoveDish(dishId)
	w.WriteHeader(http.StatusOK)

}

func setWriter(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}