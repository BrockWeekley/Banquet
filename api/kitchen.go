package kitchen

import (
	"encoding/json"
	"fmt"
	common "github.com/brockweekley/banquet"
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
	http.HandleFunc("/api", handler)
	http.HandleFunc("/api/returnMenu", returnMenu)
	http.HandleFunc("/api/serveCourse", serveCourse)
	common.PrintPositive("Kitchen is ready.")
	common.CheckForError(http.ListenAndServe(":8080", nil))
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
	readMenu(w)
}

func readMenu(w http.ResponseWriter) {
	_, fileName, _, _ := runtime.Caller(0)

	//filePath := strings.ReplaceAll(fileName, "/api/kitchen.go", "") + "/menu"
	//common.CheckForError(os.Chdir(filePath))
	//
	//files, directoryError := filepath.Glob("*")
	//common.CheckForError(directoryError)
	//files[0] = ""

	filePath := strings.ReplaceAll(fileName, "kitchen.go", "")
	common.CheckForError(os.Chdir(filePath))

	menu, fileError := os.ReadFile("menu.json")
	common.CheckForError(fileError)
	var courses []Course
	common.CheckForError(json.Unmarshal(menu, &courses))

	common.CheckForError(json.NewEncoder(w).Encode(courses))
}