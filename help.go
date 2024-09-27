package main

import (
	_ "embed"
	"os"

	"github.com/coalaura/arguments"
)

var (
	//go:embed help.txt
	helpText string
)

func help() {
	if !arguments.Bool("h", "help", false) {
		return
	}

	log.CPrint(helpText, 248, 0)

	os.Exit(0)
}
