package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

func main() {
	id, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(id)
}
