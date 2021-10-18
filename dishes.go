package main

import (
	"encoding/json"
	"os"
)

type dish struct {
	ID string
	Title string
	URL string
	ImageURLs []string
	Colors []string
	Status string
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

	return foundDish.ID == dishID
}

func serveDish(dishID string)() {

}