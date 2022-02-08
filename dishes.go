package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"encoding/json"
	"github.com/docker/distribution/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Dish struct {
	ID string
	ContainerID string
	Title string
	URL string
	ImageURLs []string
	Colors []string
	CustomStyleLocation string
	CustomTSLocation string
	IonicVariables [9]string
	Capacitor string
	Status string
	DeploymentType string
	LocalhostName string
	Token string
}

func GetDishes()(dishes []Dish) {
	file, err := os.ReadFile("./menu.json")
	CheckForError(err)
	CheckForError(json.Unmarshal(file, &dishes))

	return dishes
}

func GetDish(dishID string)(foundDish Dish) {
	file, err := os.ReadFile("./menu.json")
	CheckForError(err)
	var dishes []Dish
	CheckForError(json.Unmarshal(file, &dishes))
	for _, currentDish := range dishes {
		if currentDish.ID == dishID {
			foundDish = currentDish
		}
	}

	return foundDish
}

func CheckForExistingDishID(potentialDishID string)(status bool) {
	file, err := os.ReadFile("./menu.json")
	var dishes []Dish
	CheckForError(json.Unmarshal(file, &dishes))
	CheckForError(err)
	for _, dish := range dishes {
		if dish.ID == potentialDishID {
			return true
		}
	}
	return false
}

func AddDish(newDish Dish, existingBuild string)() {
	file, err := os.ReadFile("./menu.json")
	var dishes []Dish
	CheckForError(json.Unmarshal(file, &dishes))
	CheckForError(err)
	newDish.Title = strings.ReplaceAll(newDish.Title, " ", "_")
	dishes = append(dishes, newDish)
	dishBytes, err := json.Marshal(dishes)
	CheckForError(err)
	err = os.WriteFile("./menu.json", dishBytes, 0644)
	CheckForError(err)
	serveDish(newDish, existingBuild)
}

func RemoveDish(dishID string)(status bool) {
	file, err := os.ReadFile("./menu.json")
	CheckForError(err)
	var dishes []Dish
	var foundDish Dish
	CheckForError(json.Unmarshal(file, &dishes))
	for i, currentDish := range dishes {
		if currentDish.ID == dishID {
			foundDish = currentDish
			dishes = append(dishes[:i], dishes[i+1:]...)
		}
	}
	dishBytes, err := json.Marshal(dishes)
	CheckForError(err)
	if foundDish.ID == dishID {
		status = cleanDish(foundDish)
		if status {
			err = os.WriteFile("./menu.json", dishBytes, 0666)
			CheckForError(err)
		}
		return status
	}
	return false
}

func serveDish(dish Dish, existingBuild string)() {
	file, err := os.ReadFile("./config.json")
	CheckForError(err)
	var user user
	CheckForError(json.Unmarshal(file, &user))

	if existingBuild == "" {
		downloadRepo(dish)
	}
	generateStyling(dish, existingBuild)
	dockerize(dish, existingBuild)
	deployContainer(dish, user)
}

func cleanDish(dish Dish)(status bool) {

	if dish.DeploymentType == "firebase" {

	}
	if dish.DeploymentType == "aws" {

	}
	if dish.DeploymentType == "localhost" {
		ctx := context.Background()
		PrintPositive("Destroying Unused Docker Containers..")
		dockerClient, err := client.NewClientWithOpts(client.FromEnv)
		CheckForError(err)
		removeOptions := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks: false,
			Force: true,
		}
		imageRemoveOptions := types.ImageRemoveOptions{
			Force: true,
			PruneChildren: true,
		}
		CheckForError(dockerClient.ContainerRemove(ctx, dish.ContainerID, removeOptions))

		images, err := dockerClient.ImageList(ctx, types.ImageListOptions{All: true})
		CheckForError(err)

		for _, image := range images {
			for _, tag := range image.RepoTags {
				if tag == "banquet-" + dish.Title + ":latest" {
					_, err = dockerClient.ImageRemove(ctx, image.ID, imageRemoveOptions)
					CheckForError(err)
				}
			}
		}

		_, fileName, _, _ := runtime.Caller(0)
		filePath := strings.ReplaceAll(fileName, "dishes.go", "")
		PrintPositive("Deleting Project...")
		err = os.RemoveAll(filePath + "menu/" + dish.Title)
		if err != nil {
			PrintNegative("Unable to delete project files, this is not an issue if a sibling dish has already been deleted. Returning success.")
		}
		return true
	}
	return false
}

func downloadRepo(dish Dish) {
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

func generateStyling(dish Dish, existingBuild string) {
	workingTitle := dish.Title
	if existingBuild != "" {
		workingTitle = existingBuild
	}
	if dish.CustomStyleLocation != "" {
		original, err := os.Open(dish.CustomStyleLocation)
		CheckForError(err)
		fileName := dish.CustomStyleLocation[strings.LastIndex(dish.CustomStyleLocation, "/") + 1:]
		css, err := os.Create("./menu/" + workingTitle + "/" + workingTitle + "/src/" + fileName)
		CheckForError(err)
		_, err = io.Copy(css, original)
		CheckForError(err)
		defer CheckForError(css.Close())
	}
	if dish.CustomTSLocation != "" {
		original, err := os.Open(dish.CustomTSLocation)
		CheckForError(err)
		fileName := dish.CustomTSLocation[strings.LastIndex(dish.CustomTSLocation, "/") + 1:]
		ts, err := os.Create("./menu/" + workingTitle + "/" + workingTitle + "/src/" + fileName)
		CheckForError(err)
		_, err = io.Copy(ts, original)
		CheckForError(err)
		defer CheckForError(ts.Close())
	}

	ionic := false
	for _, variable := range dish.IonicVariables {
		if variable != "" {
			ionic = true
		}
	}
	if ionic {
		css, err := os.ReadFile("./menu/" + workingTitle + "/" + workingTitle + "/src/theme/variables.css")
		CheckForError(err)
		foundCss := string(css)
		variables := [9]string{
			"ion-color-primary",
			"ion-color-secondary",
			"ion-color-tertiary",
			"ion-color-success",
			"ion-color-warning",
			"ion-color-danger",
			"ion-color-dark",
			"ion-color-medium",
			"ion-color-light",
		}
		for i, variable := range dish.IonicVariables {
			if variable != "" {
				ionVariable := variables[i]
				r, err := regexp.Compile("--" + ionVariable + ":.+;")
				CheckForError(err)
				foundCss = r.ReplaceAllString(foundCss, "--" + ionVariable + ": " + variable + ";")
			}
		}
		err = os.WriteFile("./menu/" + workingTitle + "/" + workingTitle + "/src/theme/variables.css", []byte(foundCss), 0777)
		CheckForError(err)
	}

	if len(dish.Colors) + len(dish.ImageURLs) > 0 {
		css, err := os.Create("./menu/" + workingTitle + "/" + workingTitle + "/src/banquet.css")
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
		defer CheckForError(css.Close())
	}
}

func dockerize(dish Dish, existingBuild string) {
	workingTitle := dish.Title
	if existingBuild != "" {
		workingTitle = existingBuild
	}

	_, fileName, _, _ := runtime.Caller(0)
	filePath := strings.ReplaceAll(fileName, "dishes.go", "")
	CheckForError(os.Chdir(filePath + "menu/" + workingTitle + "/" + workingTitle + "/"))

	if existingBuild == "" {
		PrintPositive("Installing packages for project... This will probably take a while")
		cmd := exec.Command("npm", "install")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		CheckForError(cmd.Run())
	}

	PrintPositive("Building Project for Docker")
	cmd := exec.Command("npm", "run", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	CheckForError(cmd.Run())

	if dish.Capacitor != "" {
		buildMobile(dish, filePath, existingBuild)
	}

	CheckForError(os.Chdir(filePath))

	if existingBuild == "" {
		CheckForError(os.Mkdir("./menu/" + workingTitle + "/" + workingTitle + "/nginx", 0644))
	}

	nginx, err := os.Create("./menu/" + workingTitle + "/" + workingTitle + "/nginx/nginx.conf")
	CheckForError(err)
	dockerFile, err := os.Create("./menu/" + workingTitle + "/" + workingTitle + "/Dockerfile.prod")
	CheckForError(err)
	dockerIgnore, err := os.Create("./menu/" + workingTitle + "/" + workingTitle + "/.dockerignore")
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

	items, err := ioutil.ReadDir("./menu/" + workingTitle + "/" + workingTitle + "/build")
	CheckForError(err)

	buffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buffer)

	tarFiles(items, dish, tarWriter, buffer, "", existingBuild)

	dockerFile, err = os.Open("./menu/" + workingTitle + "/" + workingTitle + "/Dockerfile.prod")
	CheckForError(err)
	readDockerFile, err := ioutil.ReadAll(dockerFile)
	CheckForError(err)

	nginxFile, err := os.Open("./menu/" + workingTitle + "/" + workingTitle + "/nginx/nginx.conf")
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

func deployContainer(dish Dish, user user) {
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	CheckForError(err)
	//appImageInspect, appImage, err := dockerClient.ImageInspectWithRaw(ctx, imageID)
	if user.DeploymentType == "firebase" {
		// TODO: TBD
	}
	if user.DeploymentType == "aws" {
		// TODO: TBD
	}
	if user.DeploymentType == "localhost" {
		PrintPositive("Running container on port " + dish.LocalhostName + "...")
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

		file, err := os.ReadFile("./menu.json")
		CheckForError(err)
		var dishes []Dish
		CheckForError(json.Unmarshal(file, &dishes))
		for i, currentDish := range dishes {
			if currentDish.ID == dish.ID {
				dishes[i].ContainerID = newContainer.ID
			}
		}
		dishBytes, err := json.Marshal(dishes)
		CheckForError(err)
		_, err = os.Create("./menu.json")
		CheckForError(err)
		err = os.WriteFile("./menu.json", dishBytes, 0666)
		CheckForError(err)

		CheckForError(dockerClient.ContainerStart(ctx, newContainer.ID, types.ContainerStartOptions{}))
		pruneFilters := filters.NewArgs(
			filters.Arg("dangling", "true"))
		_, err = dockerClient.ImagesPrune(ctx, pruneFilters)
		CheckForError(err)
	}
}

func buildMobile(dish Dish, filePath string, existingBuild string) {
	workingTitle := dish.Title
	if existingBuild != "" {
		workingTitle = existingBuild
	}

	if existingBuild == "" {
		cmd := exec.Command("npm", "install", "@capacitor/core")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		CheckForError(cmd.Run())
		cmd = exec.Command("npm", "install", "@capacitor/cli")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		CheckForError(cmd.Run())
		PrintPositive("Building your application for android...")
		cmd = exec.Command("npm", "install", "@capacitor/android")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		CheckForError(cmd.Run())

		cmd = exec.Command("npx", "cap", "init", dish.Title, dish.ID, "--web-dir", "build")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		stdin, err := cmd.StdinPipe()
		CheckForError(err)
		CheckForError(cmd.Start())
		// TODO: Race condition - There has to be a better way to do this:
		time.Sleep(time.Second)
		_, err = io.WriteString(stdin, "\n")
		CheckForError(stdin.Close())
		CheckForError(cmd.Wait())
		cmd = exec.Command("npx", "cap", "add", "android")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		CheckForError(cmd.Run())
	}

	cmd := exec.Command("npx", "cap", "sync", "android")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	CheckForError(cmd.Run())
	PrintPositive("Application built successfully")

	properties, err := os.Create(filePath + "menu/" + workingTitle + "/" + workingTitle + "/android/local.properties")
	CheckForError(err)
	_, err = properties.WriteString("sdk.dir=" + dish.Capacitor)
	CheckForError(err)
	defer CheckForError(properties.Close())
	PrintPositive("Wrote SDK path to Android project. Your app is hot and ready!")
	//PrintPositive("Building your application for ios...")
}

func tarFiles(files []fs.FileInfo, dish Dish, writer *tar.Writer, buffer *bytes.Buffer, additionalPath string, existingBuild string) {
	workingTitle := dish.Title
	if existingBuild != "" {
		workingTitle = existingBuild
	}
	for _, file := range files {
		if file.IsDir() {
			items, err := ioutil.ReadDir("./menu/" + workingTitle + "/" + workingTitle + "/build/" + additionalPath + file.Name())
			CheckForError(err)
			tarFiles(items, dish, writer, buffer, additionalPath + file.Name() + "/", existingBuild)
		} else {
			currentFile, err := os.Open("./menu/" + workingTitle + "/" + workingTitle + "/build/" + additionalPath + file.Name())
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