package main

import (
	"context"
	"encoding/json"
	"fmt"
	firebase "google.golang.org/api/firebase/v1beta1"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

//context := context.Background()
//firebaseService, err := firebase.NewService(context)


type Course struct {
	Name string
	Status string
	Port int
}

type CourseRequest struct {
	ProjectName string
	GitHubURL string
	ProjectType string
	ParentID string
}

func StartServer() {
	apiServerMux := http.NewServeMux()
	apiServer := http.Server {
		Addr: fmt.Sprintf(":%v", 8080),
		Handler: apiServerMux,
	}
	apiServerMux.HandleFunc("/api", handler)
	apiServerMux.HandleFunc("/api/returnMenu", returnMenu)
	apiServerMux.HandleFunc("/api/returnAccounts", returnAccounts)
	apiServerMux.HandleFunc("/api/prepareCourse", prepareCourse)
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

func returnAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	readAccounts(w)
}

func prepareCourse(w http.ResponseWriter, r *http.Request) {
	var body CourseRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	CheckForError(err)
	CreateProject(body.ProjectName)
	directory, _ := os.Getwd()
	CloneRepository(body.GitHubURL, directory)
	InitializeStatus(body.ProjectName)
	if body.ProjectType != "local" {
		web, _, _ := createFirebaseProject(body.ProjectType, body.ParentID)
		CheckForError(json.NewEncoder(w).Encode(web))
	}
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

func readAccounts(w http.ResponseWriter) {
	firebaseService := GetService()
	projects, projectsError := firebaseService.Projects.List().Do()
	CheckForError(projectsError)
	CheckForError(json.NewEncoder(w).Encode(projects))
}

func createFirebaseProject(projectType string, parentID string) (*firebase.ProjectsWebAppsCreateCall, *firebase.ProjectsAndroidAppsCreateCall, *firebase.ProjectsIosAppsCreateCall) {
	firebaseService := GetService()
	if projectType == "web" {
		result := firebaseService.Projects.WebApps.Create(parentID, nil)
		return result, nil, nil
	}
	if projectType == "android" {
		result := firebaseService.Projects.AndroidApps.Create(parentID, nil)
		return nil, result, nil
	}
	if projectType == "ios" {
		result := firebaseService.Projects.IosApps.Create(parentID, nil)
		return nil, nil, result
	}
	return nil, nil, nil
}

func GetService() *firebase.Service {
	ctx := context.Background()
	firebaseService, err := firebase.NewService(ctx, option.WithCredentialsFile("serviceAccountKey.json"))
	CheckForError(err)
	return firebaseService
}