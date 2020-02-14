package app

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/carlmjohnson/flagext"
	"github.com/peterbourgon/ff"
)

const AppName = "go-cli"

func CLI(args []string) error {
	var app appEnv
	err := app.ParseArgs(args)
	if err == nil {
		err = app.Exec()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return err
}

func (app *appEnv) ParseArgs(args []string) error {
	fl := flag.NewFlagSet(AppName, flag.ContinueOnError)
	src := flagext.FileOrURL(flagext.StdIO, nil)
	app.src = src
	fl.Var(src, "src", "source file or URL")
	app.Logger = log.New(nil, AppName+" ", log.LstdFlags)
	fl.Var(
		flagext.Logger(app.Logger, flagext.LogVerbose),
		"verbose",
		`log debug output`,
	)

	fl.Usage = func() {
		fmt.Fprintf(fl.Output(), `go-cli - a Go CLI application template cat clone

Usage:

	go-cli [options]

Options:
`)
		fl.PrintDefaults()
		fmt.Fprintln(fl.Output(), "")
	}
	if err := ff.Parse(fl, args, ff.WithEnvVarPrefix("GO_CLI")); err != nil {
		return err
	}
	return nil
}

type appEnv struct {
	src io.ReadCloser
	*log.Logger
}

func (app *appEnv) Exec() (err error) {
	app.Println("starting")
	defer func() { app.Println("done") }()

	n, err := io.Copy(os.Stdout, app.src)
	defer func() {
		e2 := app.src.Close()
		if err == nil {
			err = e2
		}
	}()
	app.Printf("copied %d bytes\n", n)

	return err
}
