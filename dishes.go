package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	if dish.DeploymentType == "firebase" {

	}
	if dish.DeploymentType == "aws" {

	}
	if dish.DeploymentType == "localhost" {

	}
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
		client := http.Client{}
		request, err := http.NewRequest("GET", dish.URL, nil)
		CheckForError(err)
		request.Header.Add("Authorization", "token " + dish.Token)
		response, err := client.Do(request)
		CheckForError(err)
		if response.StatusCode != http.StatusOK {
			PrintNegative("Bad Status for provided GitHub URL: " + response.Status)
		}
		_, err = io.Copy(create, response.Body)

		defer CheckForError(response.Body.Close())
	} else {
		client := http.Client{}
		request, err := http.NewRequest("GET", dish.URL, nil)
		CheckForError(err)
		response, err := client.Do(request)
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
	dockerFile, err := os.Create("./menu/" + dish.Title + "/" + dish.Title + "/Dockerfile")
	dockerIgnore, err := os.Create("./menu/" + dish.Title + "/" + dish.Title + "/.dockerignore")
	CheckForError(err)
	_, err = dockerFile.WriteString("FROM node:latest\n\n" +
								"WORKDIR /app\n\n" +
								"COPY ./package.json ./\n\n" +
								"RUN npm install\n\n" +
								"COPY . .\n\n" +
								"RUN npm run build")
	CheckForError(err)
	_, err = dockerIgnore.WriteString("node_modules\n" +
										"build\n" +
										".dockerignore\n" +
										"Dockerfile\n" +
										"Dockerfile.prod")
	CheckForError(err)

	CheckForError(dockerFile.Close())

	ctx := context.Background()
	PrintPositive("Starting Docker Daemon..")
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)

	//dockerFile, err = os.Open("./menu/" + dish.Title + "/" + dish.Title + "/Dockerfile")
	//CheckForError(err)
	//readDockerFile, err := ioutil.ReadAll(dockerFile)
	//CheckForError(err)
	//
	//packageJSON, err := os.Open("./menu/" + dish.Title + "/" + dish.Title + "/package.json")
	//CheckForError(err)
	//packageJSONFile, err := ioutil.ReadAll(packageJSON)
	//CheckForError(err)

	buffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buffer)

	src := "./menu/" + dish.Title + "/" + dish.Title + "/"

	CheckForError(filepath.Walk(src, func(file string, info os.FileInfo, err error) error {
		tarHeader := &tar.Header{
			Name: info.Name(),
			Size: info.Size(),
		}
		CheckForError(tarWriter.WriteHeader(tarHeader))
		if !info.IsDir() {
			data, err := os.Open(file)
			CheckForError(err)
			_, err = io.Copy(tarWriter, data)
			CheckForError(err)
		}
		return nil
	}))
	//tarHeader := &tar.Header{
	//	Name: "Dockerfile",
	//	Size: int64(len(readDockerFile)),
	//}
	//err = tarWriter.WriteHeader(tarHeader)
	//CheckForError(err)
	//_, err = tarWriter.Write(readDockerFile)
	//CheckForError(err)
	//
	//tarPack := &tar.Header{
	//	Name: "package.json",
	//	Size: int64(len(packageJSONFile)),
	//}
	//err = tarWriter.WriteHeader(tarPack)
	//CheckForError(err)
	//_, err = tarWriter.Write(packageJSONFile)
	//CheckForError(err)

	dockerFileTarReader := bytes.NewReader(buffer.Bytes())

	buildOptions := types.ImageBuildOptions{
		Context: dockerFileTarReader,
		Dockerfile: "Dockerfile",
		Remove: 	true,
	}
	CheckForError(err)
	PrintPositive("Building Docker Image...")
	imageBuildResponse, err := dockerClient.ImageBuild(ctx, dockerFileTarReader, buildOptions)
	CheckForError(err)
	PrintPositive("Build Response: ")
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	CheckForError(err)

	defer CheckForError(dockerFile.Close())
	defer CheckForError(tarWriter.Close())
	defer CheckForError(imageBuildResponse.Body.Close())
	defer CheckForError(dockerIgnore.Close())
}