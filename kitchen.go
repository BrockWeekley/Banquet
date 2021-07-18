package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type Course struct {
	Name string
	Status string
	Port int
}

func StartServer() {
	apiServerMux := http.NewServeMux()
	apiServer := http.Server {
		Addr: fmt.Sprintf(":%v", 8080),
		Handler: apiServerMux,
	}
	apiServerMux.HandleFunc("/api", handler)
	apiServerMux.HandleFunc("/api/returnMenu", returnMenu)
	apiServerMux.HandleFunc("/api/serveCourse", serveCourse)
	PrintPositive("Kitchen is ready.")
	CheckForError(apiServer.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w)
}

func returnMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	readMenu(w)
}

func serveCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(reqBody))
	readMenu(w)
}

func readMenu(w http.ResponseWriter) {
	_, fileName, _, _ := runtime.Caller(0)

	//filePath := strings.ReplaceAll(fileName, "/api/kitchen.go", "") + "/menu"
	//CheckForError(os.Chdir(filePath))
	//
	//files, directoryError := filepath.Glob("*")
	//CheckForError(directoryError)
	//files[0] = ""

	filePath := strings.ReplaceAll(fileName, "kitchen.go", "")
	CheckForError(os.Chdir(filePath))

	menu, fileError := os.ReadFile("menu.json")
	CheckForError(fileError)
	var courses []Course
	CheckForError(json.Unmarshal(menu, &courses))

	CheckForError(json.NewEncoder(w).Encode(courses))
}