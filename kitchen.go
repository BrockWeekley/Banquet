package main

import (
	"context"
	"encoding/json"
	"fmt"
	firebase "google.golang.org/api/firebase/v1beta1"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//context := context.Background()
//firebaseService, err := firebase.NewService(context)


type Course struct {
	Name string
	Status string
	Port string
}

type Hosting struct {
	hosting []Host
}

type Host struct {
	target string
	public string
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
	InitializeStatus(body.ProjectName)
	if body.ProjectType != "Local" {
		PrintPositive(body.ProjectType)
		if body.ProjectType == "Web" {
			web := createFirebaseWebProject(body.ParentID)
			response, _ := json.Marshal(web)
			PrintPositive(string(response))
		}

		if body.ProjectType == "Android" {
			android := createFirebaseAndroidProject(body.ParentID)
			response, _ := json.Marshal(android)
			PrintPositive(string(response))
		}

		if body.ProjectType == "Ios" {
			ios := createFirebaseIosProject(body.ParentID)
			response, _ := json.Marshal(ios)
			PrintPositive(string(response))
		}

		CreateProject(body.ProjectName)
		directory, _ := os.Getwd()
		CloneRepository(body.GitHubURL, directory)
		buildProject(body.ProjectName)

		hostFirebaseSite(body.ProjectName)
		updateStatus(body.ProjectName)

	} else {
		CreateProject(body.ProjectName)
		directory, _ := os.Getwd()
		CloneRepository(body.GitHubURL, directory)
		buildProject(body.ProjectName)

		hostLocalSite()
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
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

func createFirebaseWebProject(parentID string) *firebase.ProjectsWebAppsCreateCall {
	firebaseService := GetService()
	result := firebaseService.Projects.WebApps.Create(parentID, nil)
	return result
}

func createFirebaseAndroidProject(parentID string) *firebase.ProjectsAndroidAppsCreateCall {
	firebaseService := GetService()
	result := firebaseService.Projects.AndroidApps.Create(parentID, nil)
	return result
}

func createFirebaseIosProject(parentID string) *firebase.ProjectsIosAppsCreateCall {
	firebaseService := GetService()
	result := firebaseService.Projects.IosApps.Create(parentID, nil)
	return result
}

func hostFirebaseSite(projectName string) {
	//config := &oauth2.Config{}
	//ctx := context.Background()
	//letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	//
	//state := make([]rune, rand.Int())
	//for i := range state {
	//	state[i] = letters[rand.Intn(len(letters))]
	//}
	//code := config.AuthCodeURL(string(state))
	//token, err := config.Exchange(ctx, code)

	initCommand := exec.Command("firebase", "init", "hosting")
	stdin, err := initCommand.StdinPipe()
	CheckForError(err)
	initCommand.Stdout = os.Stdout
	initCommand.Stderr = os.Stderr
	CheckForError(initCommand.Start())

	duration := time.Second * 5
	time.Sleep(duration)

	_, proceedError := io.WriteString(stdin, "Y\n")
	CheckForError(proceedError)
	time.Sleep(duration)

	_, projectError := io.WriteString(stdin, "\n")
	CheckForError(projectError)
	time.Sleep(duration)

	_, existingError := io.WriteString(stdin, "\n")
	CheckForError(existingError)
	time.Sleep(duration)

	_, publicDirectoryError := io.WriteString(stdin, "build\n")
	CheckForError(publicDirectoryError)
	time.Sleep(duration)

	_, spaError := io.WriteString(stdin, "y\n")
	CheckForError(spaError)
	time.Sleep(duration)

	_, gitHubError := io.WriteString(stdin, "N\n")
	CheckForError(gitHubError)
	time.Sleep(duration)

	_, indexError := io.WriteString(stdin, "N\n")
	CheckForError(indexError)
	time.Sleep(duration)
	stdin.Close()
	CheckForError(initCommand.Wait())

	randomInt := rand.Intn(99999)
	newSiteCommand := exec.Command("firebase", "hosting:sites:create", projectName + "-" + strconv.Itoa(randomInt))
	newSiteCommand.Stdout = os.Stdout
	newSiteCommand.Stderr = os.Stderr
	CheckForError(newSiteCommand.Run())

	targetCommand := exec.Command("firebase", "target:apply", "hosting", projectName, projectName + "-" + strconv.Itoa(randomInt))
	targetCommand.Stdout = os.Stdout
	targetCommand.Stderr = os.Stderr
	CheckForError(targetCommand.Run())

	firebaseJSON, fileError := os.ReadFile("firebase.json")
	CheckForError(fileError)
	var hosts Hosting
	CheckForError(json.Unmarshal(firebaseJSON, &hosts))
	var newHosts Hosting
	newHosts.hosting = append(hosts.hosting, Host{target: projectName, public: "build"})
	createHost, jsonError := json.Marshal(newHosts)
	CheckForError(jsonError)
	CheckForError(os.WriteFile("firebase.json", createHost, 0666))

	deployCommand := exec.Command("firebase", "deploy", "--only", "hosting:" + projectName)
	deployCommand.Stdout = os.Stdout
	deployCommand.Stderr = os.Stderr
	CheckForError(deployCommand.Run())
}

func hostLocalSite() {

}

func buildProject(projectName string) {
	installCommand := exec.Command("npm", "install")
	installCommand.Stdout = os.Stdout
	installCommand.Stderr = os.Stderr
	PrintPositive("Running npm install for project " + projectName + ", outputting stack:")
	CheckForError(installCommand.Run())

	buildCommand := exec.Command("npm", "run", "build")
	buildCommand.Stdout = os.Stdout
	buildCommand.Stderr = os.Stderr
	PrintPositive("Building project:")
	CheckForError(buildCommand.Run())
	PrintPositive("Build ran correctly")
}

func updateStatus(projectName string) {
	_, fileName, _, _ := runtime.Caller(0)
	filePath := strings.ReplaceAll(fileName, "/menu/" + projectName, "")
	CheckForError(os.Chdir(filePath))
	menu, fileError := os.ReadFile("menu.json")
	CheckForError(fileError)
	var courses []Course
	CheckForError(json.Unmarshal(menu, &courses))

	for i, course := range courses {
		if course.Name == projectName {
			courses[i].Status = "Serving"
			courses[i].Port = "https://" + projectName + ".web.app"
		}
	}

	recreate, jsonError := json.Marshal(courses)
	CheckForError(jsonError)
	CheckForError(os.WriteFile("menu.json", recreate, 0666))
}

func GetService() *firebase.Service {
	ctx := context.Background()
	firebaseService, err := firebase.NewService(ctx, option.WithCredentialsFile("serviceAccountKey.json"))
	CheckForError(err)
	return firebaseService
}