package main

import (
	"flag"
	"os"
	"testing"
)

func TestParseCliArguments(t *testing.T) {
	// Save original command-line arguments and defer their restoration
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Simulate command-line arguments
	os.Args = []string{"cmd", "-source", "/some/source", "-target", "/some/target"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Reset flags

	// Run the function and recover from os.Exit
	exitCode := runTestWithoutExit(t, func() {
		parseCliArguments() // You can call main() here if needed.
	})

	if exitCode != 0 {
		t.Errorf("Unexpected exit code: %d", exitCode)
	}
}

func runTestWithoutExit(t *testing.T, f func()) (exitCode int) {
	// Capture calls to os.Exit
	defer func() {
		if r := recover(); r != nil {
			if code, ok := r.(int); ok {
				exitCode = code
			} else {
				t.Errorf("Unexpected panic: %v", r)
			}
		}
	}()

	// Run the function that might call os.Exit()
	f()

	return 0 // No exit, return 0 as the default "exit code"
}
