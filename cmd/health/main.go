package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/naeimc/health"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s: no command\n", health.PROGRAM)
		return
	}

	switch os.Args[1] {
	case "version":
		fmt.Printf("%s v%s %s/%s\n", health.PROGRAM, health.VERSION, runtime.GOOS, runtime.GOARCH)
	case "license":
		fmt.Printf("%s\n", health.LICENSE)
	default:
		fmt.Printf("%s %s: unknown command", health.PROGRAM, os.Args[1])
	}
}
