package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"syscall"
)

func usage() {
	fmt.Fprintln(os.Stderr, "envargs -- <cmd> [ARG_TEMPLATE]...")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Argument template containing environment variables names")
	fmt.Fprintln(os.Stderr, "in the format ${ENV_NAME} will be expanded")
	os.Exit(1)
}

var tpl = regexp.MustCompile(`\$\{([^}]+)\}`)

func expand(input string) string {
	return tpl.ReplaceAllStringFunc(input, func(match string) string {
		// Extract the variable name without the ${} wrapper
		envName := match[2 : len(match)-1]
		return os.Getenv(envName)
	})
}

func main() {
	if len(os.Args) < 3 || os.Args[1] != "--" {
		usage()
	}

	cmd, args := os.Args[2], os.Args[2:]

	// don't expand the command name argument
	for i := 1; i < len(args); i++ {
		args[i] = expand(args[i])
	}

	binPath, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatal(err)
	}

	if err := syscall.Exec(binPath, args, os.Environ()); err != nil {
		log.Fatal(err)
	}
}
