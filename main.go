package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func parseExt(path string) string {
	split := strings.Split(path, ".")
	return "." + split[len(split)-1]
}

func today() string {
	now := time.Now()
	return now.Format("20060102")
}

func main() {
	fmt.Println("Usage: ./grayscale-cbz --dest=/path/to/dest <FILE>")
	dest := flag.String("dest", "/tmp/", "Destination directory for grayscale images")

	flag.Parse()

	// check if flags valid
	cbz := flag.Args()[0]
	if !pathExists(cbz) {
		fmt.Printf("path to file '%s' doesn't exist\n", cbz)
		os.Exit(1)
	}

	// check if cbz is valid
	ext := filepath.Ext(cbz)
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

	// extract cbz
	r, err := zip.OpenReader(cbz)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			os.MkdirAll(*dest+"/"+f.Name, os.ModePerm)
			continue
		}

		fmt.Printf("converting image: %s\n", f.Name)
		zipped, err := f.Open()
		if err != nil {
			fmt.Printf("error reading '%s': %s\n", f.Name, err.Error())
			continue
		}

		// grayscale images
		img, _, err := image.Decode(zipped)
		if err != nil {
			fmt.Printf("error decoding image '%s': %s\n", f.Name, err.Error())
			continue
		}

		gray := image.NewGray(img.Bounds())
		for x := 0; x < img.Bounds().Max.X; x++ {
			for y := 0; y < img.Bounds().Max.Y; y++ {
				color := img.At(x, y)
				gray.Set(x, y, color)
			}
		}

		newExt := parseExt(f.Name)
		newFile := *dest + "/" + f.Name[0:len(f.Name)-len(newExt)] + "_GRAYSCALED_" + today() + newExt
		out, err := os.Create(newFile)
		if err != nil {
			fmt.Printf("error creating destination image path '%s': %s\n", newFile, err.Error())
			continue
		}
		defer out.Close()

		err = jpeg.Encode(out, gray, nil)
		if err != nil {
			fmt.Printf("error saving image at destination '%s': %s\n", newFile, err.Error())
			continue
		}
	}

}
