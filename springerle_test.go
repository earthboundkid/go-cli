package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/carlmjohnson/springerle/txtartmpl"
)

func TestProj(t *testing.T) {
	dst := t.TempDir()
	const context = `
		{
			"author": "John Doe",
			"proj_full": "Project X",
			"proj_short": "project_x",
			"repo": "github.com/john_doe/project_x",
			"pkg": "app",
			"gversion": "1.16",
			"description": "tktk"
		}
	`
	err := txtartmpl.CLI([]string{"-context", context, "-dest", dst, "go-cli.txtar"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if err := os.Chdir(dst); err != nil {
		t.Fatalf("err: %v", err)
	}
	cmd := exec.Command("./finalize.sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(output))
	}
	if err := os.Chdir("./project_x"); err != nil {
		t.Fatalf("err: %v", err)
	}
	cmd = exec.Command("go", "test", "./...")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(output))
	}
}
