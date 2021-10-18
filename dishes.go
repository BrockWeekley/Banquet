package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func addDish(newDish dish)() {
	file, err := os.ReadFile("./menu.json")
	var dishes []dish
	CheckForError(json.Unmarshal(file, &dishes))
	CheckForError(err)
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
	create, err := os.Create("./menu/" + dish.Title + ".zip")
	CheckForError(err)
	defer CheckForError(create.Close())

	response, err := http.Get(dish.URL)
	CheckForError(err)
	defer CheckForError(response.Body.Close())
	if response.StatusCode != http.StatusOK {
		PrintNegative("Bad Status for provided GitHub URL: " + response.Status)
	}

	_, err = io.Copy(create, response.Body)

	reader, err := zip.OpenReader("./menu/" + dish.Title + ".zip")
	destination := "./menu/" + dish.Title + "/"
	CheckForError(err)

	defer CheckForError(reader.Close())
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

	CheckForError(err)
}

func dockerize(dish dish) {
	create, err := os.Create("./menu/" + dish.Title + "/Dockerfile")
	CheckForError(err)
	defer CheckForError(create.Close())
}