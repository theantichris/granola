package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/theantichris/granola/cmd"
)

func main() {
	if err := fang.Execute(context.Background(), cmd.RootCmd); err != nil {
		os.Exit(1)
	}
}
