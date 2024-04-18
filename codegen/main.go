package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("must provide 1 arg")
		os.Exit(1)
	}
	serviceName := os.Args[1]

	_, err := os.Stat(serviceName)
	if !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.Mkdir(fmt.Sprintf("services/%s", serviceName), 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.Mkdir(fmt.Sprintf("services/%s/cmd", serviceName), 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.Mkdir(fmt.Sprintf("services/%s/bin", serviceName), 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
