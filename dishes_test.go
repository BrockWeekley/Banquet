package main

import (
	"fmt"
	"testing"
)

const id = "com.banquet.tester"
const testURL = "https://api.github.com/repos/BrockWeekley/banquet-tester/zipball/master"
const testPort = "8080"
const testDeployment = false
const sdkPath = "C\\:\\\\Users\\\\Brack\\\\AppData\\\\Local\\\\Android\\\\Sdk"
const googleKey = ""

func TestAddDish(t *testing.T) {
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
		ApiKey: 				googleKey,
		Token:          		"",
	}
	AddDish(testDish)
	result := CheckForExistingDishID(id)
	if !result {
		t.Errorf("Check for Test Dish ID failed. AddDish or CheckForExistingDishID has an error.")
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