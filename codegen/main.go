// This is a very simple tool for generating the inital boilerplate for a
// service.
package main

import (
	"fmt"
	"os"
	"text/template"
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

	err = makeDirs(serviceName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = makeDockerfile(serviceName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = makeGoMod(serviceName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = makeMain(serviceName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func makeDirs(serviceName string) error {
	if err := os.Mkdir(fmt.Sprintf("services/%s", serviceName), 0755); err != nil {
		return err
	}

	if err := os.Mkdir(fmt.Sprintf("services/%s/cmd", serviceName), 0755); err != nil {
		return err
	}

	if err := os.Mkdir(fmt.Sprintf("services/%s/bin", serviceName), 0755); err != nil {
		return err
	}

	return nil
}

func makeDockerfile(serviceName string) error {
	t, err := template.ParseFiles("codegen/Dockerfile.template")
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("services/%s/Dockerfile", serviceName))
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, serviceName)
}

func makeGoMod(serviceName string) error {
	t, err := template.ParseFiles("codegen/go.mod.template")
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("services/%s/go.mod", serviceName))
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, serviceName)
}

func makeMain(serviceName string) error {
	t, err := template.ParseFiles("codegen/main.go.template")
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("services/%s/cmd/main.go", serviceName))
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, serviceName)
}
