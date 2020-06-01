package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/stoewer/go-strcase"
)

func main() {
	os.Exit(run())
}

func run() int {
	type value struct{ label, value, replacement string }
	var (
		name = value{"Project name", "go-cli", ""}
		url  = value{"Project URL", "github.com/carlmjohnson/go-cli", ""}
		err  error
	)

	for _, val := range []*value{&name, &url} {
		prompt := promptui.Prompt{
			Label: val.label,
		}

		val.replacement, err = prompt.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Aborted: %v\n", err)
			return 2
		}
		fmt.Println("ok", val.replacement)
	}
	fmt.Printf("You entered --\nProject name: %q\nProject URL: %q\n\n",
		name.replacement, url.replacement,
	)
	prompt := promptui.Prompt{
		Label:     "Initialize project?",
		Default:   "y",
		IsConfirm: true,
	}
	_, err = prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Canceled: %v\n", err)
		return 2
	}
	r := strings.NewReplacer(
		url.value, url.replacement,
		name.value, name.replacement,
		strcase.UpperSnakeCase(name.value), strcase.UpperSnakeCase(name.replacement),
	)
	rval := 0
	fnames := []string{"../main.go", "../go.mod", "../app/app.go", "readme-template.md"}
	for _, fname := range fnames {
		b, err := ioutil.ReadFile(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read %s: %v", fname, err)
			rval = 1
			continue
		}
		data := r.Replace(string(b))
		if err = ioutil.WriteFile(fname, []byte(data), os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "could not write %s: %v", fname, err)
			rval = 1
			continue
		}
	}
	if rval != 0 {
		return rval
	}
	b, err := ioutil.ReadFile("readme-template.md")
	if err = ioutil.WriteFile("../README.md", b, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "could not copy readme-template.md: %v", err)
		return 1
	}

	prompt = promptui.Prompt{
		Label:     "Remove initproj dir",
		Default:   "y",
		IsConfirm: true,
	}
	_, err = prompt.Run()
	if err != nil {
		return 0
	}

	if err = os.RemoveAll("."); err != nil {
		fmt.Fprintf(os.Stderr, "could not remove initproj: %v", err)
		return 1
	}
	return 0
}
