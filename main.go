package main

import (
	"os"

	"github.com/carlmjohnson/exitcode"
	"github.com/carlmjohnson/go-cli/app"
)

func main() {
	exitcode.Exit(app.CLI(os.Args[1:]))
}
