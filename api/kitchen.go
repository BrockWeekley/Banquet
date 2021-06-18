package kitchen

import (
	"encoding/json"
	"fmt"
	common "github.com/brockweekley/banquet"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func StartServer() {
	http.HandleFunc("/api", handler)
	http.HandleFunc("/api/returnMenu", returnMenu)
	common.PrintPositive("Kitchen is ready.")
	common.CheckForError(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w)
}

func returnMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, fileName, _, _ := runtime.Caller(0)
	filePath := strings.ReplaceAll(fileName, "/api/kitchen.go", "") + "/menu"
	common.CheckForError(os.Chdir(filePath))

	files, directoryError := filepath.Glob("*")
	common.CheckForError(directoryError)
	files[0] = ""
	common.CheckForError(json.NewEncoder(w).Encode(files))
}
