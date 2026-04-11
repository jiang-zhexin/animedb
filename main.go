package main

import (
	"fmt"
	"os"

	"github.com/jiang-zhexin/animedb/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("Exit:", err.Error())
		os.Exit(1)
	}
}
