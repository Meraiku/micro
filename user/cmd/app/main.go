package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	fmt.Println("Init from User")
}
