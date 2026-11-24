package main

import (
	"fmt"
	"os"

	"pos/internal/composition"
)

func main() {
	container := composition.NewContainer(os.Stdin, os.Stdout)

	if err := container.Runner.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Application error:", err)
		os.Exit(1)
	}
}
