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
	a, err := parseArgs(args)
	if err != nil {
		return err
	}
	if err := a.exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return err
	}

	return nil
}

func parseArgs(args []string) (*app, error) {
	fl := flag.NewFlagSet(AppName, flag.ContinueOnError)
	src := flagext.FileOrURL(flagext.StdIO, nil)
	fl.Var(src, "src", "source file or URL")
	l := log.New(nil, AppName+" ", log.LstdFlags)
	fl.Var(
		flagext.Logger(l, flagext.LogVerbose),
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
	}
	if err := ff.Parse(fl, args, ff.WithEnvVarPrefix("GO_CLI")); err != nil {
		return nil, err
	}
	a := app{src, l}
	return &a, nil
}

type app struct {
	src io.ReadCloser
	*log.Logger
}

func (a *app) exec() (err error) {
	a.Println("starting")
	defer func() { a.Println("done") }()

	n, err := io.Copy(os.Stdout, a.src)
	defer func() {
		e2 := a.src.Close()
		if err == nil {
			err = e2
		}
	}()
	a.Printf("copied %d bytes\n", n)

	return err
}
