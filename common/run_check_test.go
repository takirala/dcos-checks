package common

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func newFakeCheck(stdout string, code int, e error) *fakeCheck {
	return &fakeCheck{
		stdout: stdout,
		code:   code,
		e:      e,
	}
}

type fakeCheck struct {
	stdout string
	code   int
	e      error
}

func (f fakeCheck) ID() string {
	return "fakeCheck"
}

func (f fakeCheck) Run(context.Context, *CLIConfigFlags) (string, int, error) {
	return f.stdout, f.code, f.e
}

// taken from https://talks.golang.org/2014/testing.slide#23
func TestRunCheckFail(t *testing.T) {
	errMsg := "some error text"
	f := newFakeCheck("", 2, errors.New(errMsg))

	if os.Getenv("BE_CRASHER") == "1" {
		RunCheck(context.TODO(), f)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunCheckFail")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	expectedErrorMsg := fmt.Sprintf("Error executing fakeCheck: %s\n", errMsg)
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		if stdout.String() != "" {
			t.Fatalf("expect empty stdout. Got %s", stdout.String())
		}

		if stderr.String() != expectedErrorMsg {
			t.Fatalf("expect \"%s\". Got \"%s\"", expectedErrorMsg, stderr.String())
		}
		return
	}

	t.Fatalf("expect exit code 2. Got error %s", err)
}

func TestRunCheckSuccess(t *testing.T) {
	expectedStdout := "all is good"
	f := newFakeCheck(expectedStdout, 0, nil)

	if os.Getenv("BE_CRASHER") == "1" {
		RunCheck(context.TODO(), f)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunCheckSuccess")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		t.Fatalf("expect exit code 0. Got %s. Stdout %s, stderr %s", e.Error(), stdout.String(), stderr.String())
	}

	if strings.Trim(stdout.String(), "\n") != expectedStdout {
		t.Fatalf("expect stdout: \"%s\". Got \"%s\"", expectedStdout, stdout.String())
	}

	if stderr.String() != "" {
		t.Fatalf("Stderr must be empty. Got \"%s\"", stderr.String())
	}
}
