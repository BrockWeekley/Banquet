package main

import (
	"fmt"
	"os"
	"testing"
)

const id = "com.banquet.tester"
const testURL = "https://api.github.com/repos/BrockWeekley/banquet-tester/zipball/master"
const testPort = "8080"
const sdkPath = ""
// sdkPath example: "C\\:\\\\Users\\\\User\\\\AppData\\\\Local\\\\Android\\\\Sdk"

func TestAddDish(t *testing.T) {

	//TODO: For some reason, running AddDish sequentially fails with busy resources error under npm run build

	if id == "com.banquet.tester" {
		testDish := Dish{
			ID:             		id,
			Title:          		"test dish",
			URL:            		testURL,
			ImageURLs:      		[]string{""},
			Colors:         		[]string{""},
			CustomStyleLocation: 	"",
			CustomTSLocation: 		"",
			IonicVariables: 		[9]string{"#FFFFFF", "", "", "", "", "", "", "", ""},
			Capacitor: 				sdkPath,
			Status:         		"stopped",
			DeploymentType: 		"localhost",
			LocalhostName:  		testPort,
			Token:          		"",
		}

		AddDish(testDish, "")
		result := CheckForExistingDishID(id)
		if !result {
			t.Errorf("Check for Test Dish ID failed. AddDish or CheckForExistingDishID has an error.")
		}
	} else if id == "com.banquet.tester2" {
		err := os.WriteFile("./banquet.ts", []byte("export const title = \"Art Studio\";"), 0666)
		CheckForError(err)
		testDish1 := Dish{
			ID:             		id,
			Title:          		"test dish 2",
			URL:            		testURL,
			ImageURLs:      		[]string{
				"https://wallpapercave.com/wp/wp3779091.jpg",
				"https://previews.123rf.com/images/olegdudko/olegdudko1911/olegdudko191100003/133626462-pinsel-kunst-malen-kreativit%C3%A4t-handwerk-hintergr%C3%BCnde-ausstellung.jpg?fj=1",
			},
			Colors:         		[]string{"#F34221", "#09C0EC"},
			CustomStyleLocation: 	"",
			CustomTSLocation: 		"./banquet.ts",
			IonicVariables: 		[9]string{"#FFFFF1", "", "", "", "", "", "", "", ""},
			Capacitor: 				sdkPath,
			Status:         		"stopped",
			DeploymentType: 		"localhost",
			LocalhostName:  		"8081",
			Token:          		"",
		}

		AddDish(testDish1, "test_dish")
		result1 := CheckForExistingDishID(id)
		if !result1 {
			t.Errorf("Check for Test Dish IDs failed. AddDish or CheckForExistingDishID has an error with multiple dishes.")
		}
	} else {
		err := os.WriteFile("./banquet.ts", []byte("export const title = \"Tech Co\";"), 0666)
		CheckForError(err)
		testDish2 := Dish{
			ID:             		id,
			Title:          		"test dish 3",
			URL:            		testURL,
			ImageURLs:      		[]string{"https://mcdn.wallpapersafari.com/medium/37/3/1REU5K.jpg", "https://wallpapercave.com/wp/wp2848568.jpg"},
			Colors:         		[]string{"#0C1041", "#1EAC58"},
			CustomStyleLocation: 	"",
			CustomTSLocation: 		"./banquet.ts",
			IonicVariables: 		[9]string{"#FFFFF2", "", "", "", "", "", "", "", ""},
			Capacitor: 				sdkPath,
			Status:         		"stopped",
			DeploymentType: 		"localhost",
			LocalhostName:  		"8082",
			Token:          		"",
		}
		AddDish(testDish2, "test_dish")
		result2 := CheckForExistingDishID(id)
		if !result2 {
			t.Errorf("Check for Test Dish IDs failed. AddDish or CheckForExistingDishID has an error with multiple dishes.")
		}
	}

}

func TestGetDishes(t *testing.T) {
	dishes := GetDishes()
	if len(dishes) < 1 {
		t.Errorf("No dishes found.")
	}
	for _, dish := range dishes {
		fmt.Print(dish)
		fmt.Print("\n")
	}
}

func TestGetDish(t *testing.T) {
	dish := GetDish(id)
	if dish.ID == "" {
		t.Errorf("No Dish found.")
	}
}

func TestRemoveDish(t *testing.T) {
	result := RemoveDish(id)
	if !result {
		t.Errorf("Unable to remove Test Dish.")
	}
}
