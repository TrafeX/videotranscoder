package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/fatih/color"
)

func parseCliArguments() (string, string, bool) {
	var sourcePath string
	var targetPath string
	var overwriteExisting bool

	flag.StringVar(&sourcePath, "source", "", "Path to source folder")
	flag.StringVar(&targetPath, "target", "", "Path to target folder")
	flag.BoolVar(&overwriteExisting, "overwrite", false, "Overwrite existing files")

	flag.Parse()

	if sourcePath == "" {
		fmt.Fprintf(os.Stderr, "missing required -source argument\n")
		os.Exit(2)
	}
	if targetPath == "" {
		fmt.Fprintf(os.Stderr, "missing required -target argument\n")
		os.Exit(2)
	}

	return sourcePath, targetPath, overwriteExisting
}

func transcodeFile(sourceFile string, targetFile string) ([]byte, error) {
	cmd := exec.Command("ffmpeg", "-y", "-i", sourceFile, "-c:v", "prores_ks", "-profile:v", "3", "-qscale:v", "9", "-c:a", "pcm_s16le", targetFile)
	output, err := cmd.CombinedOutput()
	return output, err
}

func main() {

	sourcePath, targetPath, overwriteExisting := parseCliArguments()

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "target folder %s does not exist\n", targetPath)
		os.Exit(2)
	}
	targetPath = filepath.Join(targetPath, path.Base(sourcePath)+"-transcoded")

	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d videos to convert in source folder\n", len(files))
	fmt.Println("Writing transcoded files to target folder:", targetPath)

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		err = os.Mkdir(targetPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	for nr, file := range files {

		sourceFile := filepath.Join(sourcePath, file.Name())
		targetFile := filepath.Join(targetPath, file.Name())

		if _, err := os.Stat(targetFile); !os.IsNotExist(err) && !overwriteExisting {
			color.Yellow("Target file %s already exists, skipping %s..", targetFile, file.Name())
			continue
		}
		fmt.Println("Processing file", file.Name())

		output, err := transcodeFile(sourceFile, targetFile)
		if err != nil {
			log.Fatal(err, string(output))
		}

		color.Blue("Output:\n%s", string(output))
		fmt.Printf("Processed file %d of %d\n\n", nr+1, len(files))
	}
	fmt.Println("All done.")

}
