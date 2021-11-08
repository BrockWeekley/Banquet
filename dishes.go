package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type dish struct {
	ID string
	Title string
	URL string
	ImageURLs []string
	Colors []string
	Status string
	DeploymentType string
	LocalhostName string
	Token string
}

func getDishes()(dishes []dish) {
	file, err := os.ReadFile("./menu.json")
	CheckForError(err)
	CheckForError(json.Unmarshal(file, &dishes))

	return dishes
}

func getDish(dishID string)(foundDish dish) {
	file, err := os.ReadFile("./menu.json")
	CheckForError(err)
	var dishes []dish
	CheckForError(json.Unmarshal(file, &dishes))
	for _, currentDish := range dishes {
		if currentDish.ID == dishID {
			foundDish = currentDish
		}
	}

	return foundDish
}

func checkForExistingDishID(potentialDishID string)(status bool) {
	file, err := os.ReadFile("./menu.json")
	var dishes []dish
	CheckForError(json.Unmarshal(file, &dishes))
	CheckForError(err)
	for _, dish := range dishes {
		if dish.ID == potentialDishID {
			return true
		}
	}
	return false
}

func addDish(newDish dish)() {
	file, err := os.ReadFile("./menu.json")
	var dishes []dish
	CheckForError(json.Unmarshal(file, &dishes))
	CheckForError(err)
	newDish.Title = strings.ReplaceAll(newDish.Title, " ", "_")
	dishes = append(dishes, newDish)
	dishBytes, err := json.Marshal(dishes)
	CheckForError(err)
	err = os.WriteFile("./menu.json", dishBytes, 0644)
	CheckForError(err)
	serveDish(newDish)
}

func removeDish(dishID string)(status bool) {
	file, err := os.ReadFile("./menu.json")
	CheckForError(err)
	var dishes []dish
	var foundDish dish
	CheckForError(json.Unmarshal(file, &dishes))
	for i, currentDish := range dishes {
		if currentDish.ID == dishID {
			foundDish = currentDish
			dishes = append(dishes[:i], dishes[i+1:]...)
		}
	}
	dishBytes, err := json.Marshal(dishes)
	CheckForError(err)
	err = os.WriteFile("./menu.json", dishBytes, 0666)
	CheckForError(err)
	cleanDish(foundDish)
	return foundDish.ID == dishID
}

func serveDish(dish dish)() {
	file, err := os.ReadFile("./config.json")
	CheckForError(err)
	var user user
	CheckForError(json.Unmarshal(file, &user))

	downloadRepo(dish)
	generateStyling(dish)
	dockerize(dish)
	deployContainer(dish, user)
}

func cleanDish(dish dish)() {
	file, err := os.ReadFile("./config.json")
	CheckForError(err)
	var user user
	CheckForError(json.Unmarshal(file, &user))

	if dish.DeploymentType == "firebase" {

	}
	if dish.DeploymentType == "aws" {

	}
	if dish.DeploymentType == "localhost" {

	}
}

func downloadRepo(dish dish) {
	if _, err := os.Stat("./menu"); os.IsNotExist(err) {
		CheckForError(os.Mkdir("./menu", os.ModeDir))
	}
	create, err := os.Create("./menu/" + dish.Title + ".zip")
	CheckForError(err)

	if dish.Token != "" {
		httpClient := http.Client{}
		request, err := http.NewRequest("GET", dish.URL, nil)
		CheckForError(err)
		request.Header.Add("Authorization", "token " + dish.Token)
		response, err := httpClient.Do(request)
		CheckForError(err)
		if response.StatusCode != http.StatusOK {
			PrintNegative("Bad Status for provided GitHub URL: " + response.Status)
		}
		_, err = io.Copy(create, response.Body)

		defer CheckForError(response.Body.Close())
	} else {
		httpClient := http.Client{}
		request, err := http.NewRequest("GET", dish.URL, nil)
		CheckForError(err)
		response, err := httpClient.Do(request)
		CheckForError(err)
		if response.StatusCode != http.StatusOK {
			PrintNegative("Bad Status for provided GitHub URL: " + response.Status)
		}
		_, err = io.Copy(create, response.Body)

		defer CheckForError(response.Body.Close())
	}

	reader, err := zip.OpenReader("./menu/" + dish.Title + ".zip")
	destination := "./menu/" + dish.Title + "/"
	CheckForError(err)

	var filenames []string

	for _, foundFile := range reader.File {

		path := filepath.Join(destination, foundFile.Name)

		filenames = append(filenames, path)

		if foundFile.FileInfo().IsDir() {
			CheckForError(os.MkdirAll(path, os.ModePerm))
			continue
		}

		CheckForError(os.MkdirAll(filepath.Dir(path), os.ModePerm))

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, foundFile.Mode())
		CheckForError(err)

		rc, err := foundFile.Open()
		CheckForError(err)

		_, err = io.Copy(outFile, rc)

		CheckForError(err)

		CheckForError(outFile.Close())
		CheckForError(rc.Close())
	}

	folders, err := os.ReadDir("./menu/" + dish.Title)
	CheckForError(err)
	err = os.Rename("./menu/" + dish.Title + "/" + folders[0].Name(), "./menu/" + dish.Title + "/" + dish.Title)
	CheckForError(err)
	defer CheckForError(reader.Close())
	defer CheckForError(create.Close())
	defer CheckForError(os.Remove("./menu/" + dish.Title + ".zip"))
}

func generateStyling(dish dish) {
	css, err := os.Create("./menu/" + dish.Title + "/" + dish.Title + "/src/banquet.css")
	CheckForError(err)
	ts, err := os.Create("./menu/" + dish.Title + "/" + dish.Title + "/src/banquet.ts")
	CheckForError(err)
	for index, image := range dish.ImageURLs {
		if image != "" {
			_, err = css.WriteString(".img" + strconv.Itoa(index + 1) + "{background-image: url(" + image + ");}\n")
			CheckForError(err)
		}
	}
	for index, color := range dish.Colors {
		if color != "" {
			_, err = css.WriteString(".color" + strconv.Itoa(index + 1) + "{color: " + color + ";}\n")
			CheckForError(err)
			_, err = css.WriteString(".bcolor" + strconv.Itoa(index + 1) + "{background-color: " + color + ";}\n")
			CheckForError(err)
		}
	}
	_, err = ts.WriteString("export const title = \"Tech Co\";")
	defer CheckForError(css.Close())
	defer CheckForError(ts.Close())
}

func dockerize(dish dish) {
	_, fileName, _, _ := runtime.Caller(0)
	filePath := strings.ReplaceAll(fileName, "dishes.go", "")
	CheckForError(os.Chdir(filePath + "menu/" + dish.Title + "/" + dish.Title + "/"))
	PrintPositive("Installing packages for project... This will probably take a while")
	cmd := exec.Command("npm", "install", "typescript", "-g")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	CheckForError(cmd.Run())
	cmd = exec.Command("npm", "install", "react-scripts@latest", "-g")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	CheckForError(cmd.Run())
	cmd = exec.Command("npm", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	CheckForError(cmd.Run())
	cmd = exec.Command("npm", "run", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	CheckForError(cmd.Run())
	CheckForError(os.Chdir(filePath))

	CheckForError(os.Mkdir("./menu/" + dish.Title + "/" + dish.Title + "/nginx", 0644))
	nginx, err := os.Create("./menu/" + dish.Title + "/" + dish.Title + "/nginx/nginx.conf")
	CheckForError(err)
	dockerFile, err := os.Create("./menu/" + dish.Title + "/" + dish.Title + "/Dockerfile.prod")
	CheckForError(err)
	dockerIgnore, err := os.Create("./menu/" + dish.Title + "/" + dish.Title + "/.dockerignore")
	CheckForError(err)
	_, err = dockerFile.WriteString("FROM node:16.13.0 as build\n\n" +
								"WORKDIR /app\n\n" +
								"COPY . ./\n\n" +
								"FROM nginx:1.17.8-alpine\n" +
								"COPY --from=build /app /usr/share/nginx/html\n" +
								"RUN chmod -R 765 /usr/share/nginx/html\n" +
								"RUN rm /etc/nginx/conf.d/default.conf\n" +
								"COPY ./nginx.conf /etc/nginx/conf.d\n" +
								"EXPOSE 80\n" +
								"CMD [\"nginx\", \"-g\", \"daemon off;\"]")
	CheckForError(err)
	_, err = dockerIgnore.WriteString("node_modules\n" +
										"build\n" +
										".dockerignore\n" +
										"Dockerfile\n" +
										"Dockerfile.prod")
	CheckForError(err)
	_, err = nginx.WriteString("server {\n\n  listen 80;\n\n  location / {\n    " +
		"root   /usr/share/nginx/html;\n    index  index.html index.htm;\n\n   " +
		"try_files $uri /index.html; \n  }\n\n  error_page   500 502 503 504  /50x.html;\n\n  " +
		"location = /50x.html {\n    root   /usr/share/nginx/html;\n  }\n\n}")
	CheckForError(err)

	ctx := context.Background()
	PrintPositive("Starting Docker Daemon..")
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)

	items, err := ioutil.ReadDir("./menu/" + dish.Title + "/" + dish.Title + "/build")
	CheckForError(err)

	buffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buffer)

	tarFiles(items, dish, tarWriter, buffer, "")

	dockerFile, err = os.Open("./menu/" + dish.Title + "/" + dish.Title + "/Dockerfile.prod")
	CheckForError(err)
	readDockerFile, err := ioutil.ReadAll(dockerFile)
	CheckForError(err)

	nginxFile, err := os.Open("./menu/" + dish.Title + "/" + dish.Title + "/nginx/nginx.conf")
	CheckForError(err)
	readNginxFile, err := ioutil.ReadAll(nginxFile)
	CheckForError(err)

	tarHeader := &tar.Header{
		Name: "Dockerfile.prod",
		Size: int64(len(readDockerFile)),
	}
	err = tarWriter.WriteHeader(tarHeader)
	CheckForError(err)
	_, err = tarWriter.Write(readDockerFile)
	CheckForError(err)

	tarNgin:= &tar.Header{
		Name: "nginx.conf",
		Size: int64(len(readNginxFile)),
	}
	err = tarWriter.WriteHeader(tarNgin)
	CheckForError(err)
	_, err = tarWriter.Write(readNginxFile)
	CheckForError(err)

	dockerFileTarReader := bytes.NewReader(buffer.Bytes())

	buildOptions := types.ImageBuildOptions{
		Context: dockerFileTarReader,
		Tags: []string{"banquet-" + dish.Title},
		Dockerfile: "Dockerfile.prod",
		Remove: 	true,
	}
	CheckForError(err)
	PrintPositive("Building Docker Image...")
	imageBuildResponse, err := dockerClient.ImageBuild(ctx, dockerFileTarReader, buildOptions)
	CheckForError(err)
	PrintPositive("Build Response: ")
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	stringRead := new(bytes.Buffer)
	_, err = stringRead.ReadFrom(imageBuildResponse.Body)
	CheckForError(err)
	//buildResponse := stringRead.String()
	//idIndex := strings.Index(buildResponse, "\"ID\":\"")
	//idIndex += 6
	//endID := strings.Index(buildResponse[idIndex:len(buildResponse) - 1], "\"")
	//imageID := buildResponse[idIndex: endID]
	//fmt.Println(imageID)
	CheckForError(err)

	defer CheckForError(dockerFile.Close())
	defer CheckForError(dockerIgnore.Close())
}

func deployContainer(dish dish, user user) {
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	CheckForError(err)
	//appImageInspect, appImage, err := dockerClient.ImageInspectWithRaw(ctx, imageID)
	if user.DeploymentType == "firebase" {
		fmt.Println("Hi")
	}
	if user.DeploymentType == "aws" {
		fmt.Println("hi")
	}
	if user.DeploymentType == "localhost" {
		PrintPositive("Running container on port 8080...")
		newContainer, err := dockerClient.ContainerCreate(ctx, &container.Config{
			Image: "banquet-" + dish.Title,
			ExposedPorts: nat.PortSet{
				"80/tcp": struct{}{},
			},
		}, &container.HostConfig{
			PortBindings: nat.PortMap{
				"80/tcp": []nat.PortBinding{
					{
						HostIP: "0.0.0.0",
						HostPort: dish.LocalhostName,
					},
				},
			},
		}, nil, nil, "")
		CheckForError(err)
		CheckForError(dockerClient.ContainerStart(ctx, newContainer.ID, types.ContainerStartOptions{}))
	}
}

func tarFiles(files []fs.FileInfo, dish dish, writer *tar.Writer, buffer *bytes.Buffer, additionalPath string) {
	for _, file := range files {
		if file.IsDir() {
			items, err := ioutil.ReadDir("./menu/" + dish.Title + "/" + dish.Title + "/build/" + additionalPath + file.Name())
			CheckForError(err)
			tarFiles(items, dish, writer, buffer, additionalPath + file.Name() + "/")
		} else {
			currentFile, err := os.Open("./menu/" + dish.Title + "/" + dish.Title + "/build/" + additionalPath + file.Name())
			CheckForError(err)
			currentFileData, err := ioutil.ReadAll(currentFile)
			CheckForError(err)

			tarHeader := &tar.Header{
				Name: additionalPath + file.Name(),
				Size: int64(len(currentFileData)),
			}
			err = writer.WriteHeader(tarHeader)
			CheckForError(err)
			_, err = writer.Write(currentFileData)
			CheckForError(err)
		}
	}
}