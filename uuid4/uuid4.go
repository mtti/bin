package main

import (
	"fmt"
	"os"

	"github.com/gofrs/uuid"
)

func main() {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(id)
}
