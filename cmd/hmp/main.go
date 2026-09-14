package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/enotinc/hmp/object"
	"github.com/enotinc/hmp/repl"
	"github.com/enotinc/hmp/runner"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		repl.Start(os.Stdin, os.Stdout)
	}

	input, err := readFrom(args)
	if err != nil {
		panic(err)
	}

	env := object.NewEnviroment()
	runner.Run(input, env, os.Stdout)
}

const hmp_ext = ".hmp"

func readFrom(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("To many arguments")
	}

	file := args[1]
	if _, err := os.Stat(file); err != nil {
		return "", fmt.Errorf("file %s is not exists", file)
	}

	ext := filepath.Ext(file)
	if ext != hmp_ext {
		return "", fmt.Errorf("unsupported file: %s", file)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
