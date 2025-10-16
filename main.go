/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/gingray/swisstools/cmd"
	"os"
)

func main() {
	if len(os.Args) == 1 && os.Getenv("SWISS_TEST") != "" {
		os.Args = append(os.Args, os.Getenv("SWISS_TEST"))
	}
	cmd.Execute()
}
