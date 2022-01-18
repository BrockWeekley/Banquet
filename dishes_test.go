package main

import (
	"fmt"
	"testing"
)

const id = "unique_test1"
const testURL = "https://api.github.com/repos/BrockWeekley/banquet-tester/zipball/master"
const testPort = "8080"
const testDeployment = false

func TestAddDish(t *testing.T) {
	testDish := Dish{
		ID:             id,
		Title:          "test dish",
		URL:            testURL,
		ImageURLs:      []string{""},
		Colors:         []string{""},
		Status:         "stopped",
		DeploymentType: "localhost",
		LocalhostName:  testPort,
		Token:          "",
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
		t.Errorf("No dish found.")
	}
}

func TestRemoveDish(t *testing.T) {
	result := RemoveDish(id)
	if !result {
		t.Errorf("Unable to remove Test Dish.")
	}
}
