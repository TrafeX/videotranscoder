package main

import (
	"os"
	"testing"
)

func TestBasicWorking(t *testing.T) {
	os.Args = []string{"./videotranscoder", "-help"}
	main()
}
