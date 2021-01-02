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
	"runtime"

	"github.com/fatih/color"
)

var (
	name       = "videotranscoder"
	build      = "none"
	version    = "dev-build"
	goVersion  = runtime.Version()
	versionStr = fmt.Sprintf("%s version %v, build %v %v", name, version, build, goVersion)
)

func parseCliArguments() (string, string, bool, bool) {
	var sourcePath string
	var targetPath string
	var overwriteExisting bool
	var verboseOutput bool
	var printVersion bool
	var printHelp bool

	flag.StringVar(&sourcePath, "source", "", "Path to source folder")
	flag.StringVar(&targetPath, "target", "", "Path to target folder")
	flag.BoolVar(&overwriteExisting, "overwrite", false, "Overwrite existing files in the target folder")
	flag.BoolVar(&verboseOutput, "verbose", false, "Show output of ffmpeg")
	flag.BoolVar(&printVersion, "version", false, "Print version information")
	flag.BoolVar(&printHelp, "help", false, "Print help and usage information")

	flag.Parse()

	if printVersion {
		fmt.Println(versionStr)
		os.Exit(0)
	}

	if printHelp {
		fmt.Printf("%s - transcode media to Apple's ProRes encoding\n", name)
		fmt.Printf("usage: %s -source /path/to/media/files -target /path/to/output/folder\n\n", name)
		fmt.Println("options:")
		flag.PrintDefaults()
		fmt.Println(versionStr)
		os.Exit(0)
	}

	if sourcePath == "" {
		fmt.Fprintf(os.Stderr, "missing required -source argument\ntry '%s -help' for usage information\n", name)
		os.Exit(2)
	}
	if targetPath == "" {
		fmt.Fprintf(os.Stderr, "missing required -target argument\ntry '%s -help' for usage information\n", name)
		os.Exit(2)
	}

	return sourcePath, targetPath, overwriteExisting, verboseOutput
}

func transcodeFile(sourceFile string, targetFile string) ([]byte, error) {
	cmd := exec.Command("ffmpeg", "-y", "-i", sourceFile, "-c:v", "prores_ks", "-profile:v", "3", "-qscale:v", "9", "-c:a", "pcm_s16le", targetFile)
	output, err := cmd.CombinedOutput()
	return output, err
}

func main() {

	sourcePath, targetPath, overwriteExisting, verboseOutput := parseCliArguments()

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		log.Fatalf("target folder %s does not exist", targetPath)
	}
	targetPath = filepath.Join(targetPath, path.Base(sourcePath)+"-transcoded")

	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d videos to convert in source folder", len(files))
	log.Printf("Writing transcoded files to target folder: %s", targetPath)

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
			log.Printf("Target file %s already exists, skipping %s..", targetFile, file.Name())
			continue
		}
		log.Println("Processing file", file.Name())

		output, err := transcodeFile(sourceFile, targetFile)
		if err != nil {
			log.Fatal(err, string(output))
		}

		if verboseOutput {
			color.Blue("Result:\n%s", string(output))
		}
		log.Printf("Done processing file %s (%d of %d)", file.Name(), nr+1, len(files))
	}
	log.Println("All done.")
}
