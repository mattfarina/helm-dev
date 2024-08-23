package main

import (
	"fmt"
	"os"

	"github.com/mattfarina/helm-dev/cmd/helm-dev/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
