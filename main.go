package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// The execution steps of this program are to:
// 1) take a CBR file as input and extract its contents
// 2) iterate through each file and turn into a grayscale image
// 3) save new content into target directory

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func main() {
	fmt.Println("Usage: ./grayscale-cbr --dest=/path/to/dest <FILE>")
	dest := flag.String("dest", "/tmp/", "Destination directory for grayscale images")

	flag.Parse()

	// check if flags valid
	cbr := flag.Args()[0]
	if !pathExists(cbr) {
		fmt.Printf("path to file '%s' doesn't exist\n", cbr)
		os.Exit(1)
	}

	// check if cbr is valid cbr
	ext := filepath.Ext(cbr)
	if ext != ".cbz" {
		fmt.Printf("invalid file extension: '%s'\n", ext)
		os.Exit(1)
	}

	// check if dest exists (if no create or exit with error)
	if !pathExists(*dest) {
		fmt.Printf("destination folder '%s' doesn't exist. Creating...", *dest)
		err := os.MkdirAll(*dest, os.ModePerm)
		if err != nil {
			fmt.Printf("error creating directory: %s", err.Error())
			os.Exit(1)
		}
		fmt.Println("done!")
	}

	// rename to zip
	// zipName := cbr[0:len(cbr)-len(ext)] + ".zip"
	// fmt.Println(zipName)
	// err := os.Rename(cbr, zipName)
	// if err != nil {
	// 	fmt.Printf("unable to convert cbr to zip: '%s'\n", err.Error())
	// 	os.Exit(1)
	// }

	// extract cbr
	r, err := zip.OpenReader(cbr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()

	for _, f := range r.File {
		fmt.Println(f.Name)
	}

	// grayscale images
	// save to destination -> "original_name_GRAYSCALED20201001"

}
