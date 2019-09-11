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

func CLI(args []string) error {
	fl := flag.NewFlagSet("app", flag.ContinueOnError)
	src := flagext.FileOrURL(flagext.StdIO, nil)
	fl.Var(src, "src", "source file or URL")
	verbose := fl.Bool("verbose", false, "log debug output")
	fl.Usage = func() {
		fmt.Fprintf(fl.Output(), `go-cli - a Go CLI application template cat clone

Usage:

	go-cli [options]

Options:
`)
		fl.PrintDefaults()
	}
	if err := ff.Parse(fl, args, ff.WithEnvVarPrefix("GO_CLI")); err != nil {
		return err
	}

	return appExec(src, *verbose)
}

func appExec(src io.ReadCloser, verbose bool) error {
	l := nooplogger
	if verbose {
		l = log.New(os.Stderr, "go-cli", log.LstdFlags).Printf
	}
	a := app{src, l}
	if err := a.exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return err
	}
	return nil
}

type logger = func(format string, v ...interface{})

func nooplogger(format string, v ...interface{}) {}

type app struct {
	src io.ReadCloser
	log logger
}

func (a *app) exec() (err error) {
	a.log("starting")
	defer func() { a.log("done") }()

	n, err := io.Copy(os.Stdout, a.src)
	defer func() {
		e2 := a.src.Close()
		if err == nil {
			err = e2
		}
	}()
	a.log("copied %d bytes", n)

	return err
}
