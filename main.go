package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"

	"golang_test_task1/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}
