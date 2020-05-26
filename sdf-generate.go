package main

import (
	"flag"
	"fmt"
	"github.com/maxfish/SDF-Generator/sdf"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

var version = "1.0.0"

var inputPath, outputPath, channels string
var spread, threshold float64
var downscale int
var canOverwrite bool

var channelFlags [4]bool

func main() {
	fmt.Printf("=== SDF Generator v%s (github.com/maxfish/sdf-generator) ===\n", version)

	flag.StringVar(&inputPath, "input", "", "Specify an input filename or a folder path.")
	flag.StringVar(&outputPath, "output", "", "Specify an output filename or a folder path. If input is a path then output must be a path.")
	flag.Float64Var(&spread, "spread", 4.0, "Specify the spread of the distance field. The spread is the maximum distance in pixels that will be scanned looking for a nearby edge.")
	flag.IntVar(&downscale, "downscale", 1, "Sets the factor by which to downscale the image during processing. The output image will be smaller than the input image by this factor, rounded downwards.\nNote: For greater accuracy, images to be used as input for a distance field are often generated at higher resolution.")
	flag.StringVar(&channels, "channels", "A", "Specify which channels of the input image can contribute to defining the \"inner\" area of the shape.\nAccepted values are R,G,B,A and they can be specified separated with a comma.\nE.g. 'R,A' means that the algorithm will consider a pixel \"inside\" the shape if the Red channel, or the Alpha channel, are above the threshold. ")
	flag.Float64Var(&threshold, "threshold", 0.5, "Specify the threshold applied to the channels for one pixel to be considered \"inside\" the source shape.\nThe accepted values go from 0.0 to 1.0.")
	flag.BoolVar(&canOverwrite, "overwrite", false, "Specify if the output file, when it already exists, can be overwritten.\nWARNING: this flag is applied to the whole operation and it can delete many pre-existing images.")
	flag.Parse()

	if inputPath == "" {
		exitWithError("The 'input' parameter cannot be empty.", nil)
	}
	if spread <= 0 {
		exitWithError("The 'spread' parameter has to be a positive number greater than zero", nil)
	}
	if downscale <= 0 {
		exitWithError("The 'downscale' parameter has to be a positive number greater than zero", nil)
	}
	if threshold < 0 || threshold > 1.0 {
		exitWithError("The 'threshold' parameter has to be a positive number within 0.0 and 1.0", nil)
	}

	ch := strings.Split(strings.ToLower(channels), ",")
	for _, c := range ch {
		if c == "r" {
			channelFlags[0] = true
		} else if c == "g" {
			channelFlags[1] = true
		} else if c == "b" {
			channelFlags[2] = true
		} else if c == "a" {
			channelFlags[3] = true
		} else {
			exitWithError("'source-ch' can only contain one, or more, of these values separated by a comma: R,G,B,A", nil)
		}
	}

	if pathIsFolder(inputPath) {
		if !pathIsFolder(outputPath) {
			exitWithError("When 'input' is a folder then 'output' must be a folder as well.", nil)
		}
		convertFolder()
	} else {
		if pathIsFolder(outputPath) {
			exitWithError("When 'input' is a folder then 'output' has to be a folder as well.", nil)
		}
		convertFile(inputPath, outputPath)
	}
}

func convertFolder() {
	err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".png" {
			_, file := filepath.Split(path)
			convertFile(path, filepath.Join(outputPath, file))
		}
		return nil
	})
	if err != nil {
		exitWithError("while iterating the files within the folder", nil)
	}
}

func convertFile(inputFilename string, outputFilename string) {
	flag.Parse()
	fmt.Printf("Processing file '%s'...", inputFilename)

	// Checks if we can overwrite the file
	if pathExists(outputFilename) && !canOverwrite {
		exitWithError("Destination file already exists. If you want to allow it to be overwritten then specify the 'overwrite' parameter.", nil)
	}

	// load input file
	inputFile, err := os.Open(inputFilename)
	if err != nil {
		panic(err)
	}

	// decode input file
	inputImage, _, err := image.Decode(inputFile)
	if err != nil {
		exitWithError("While decoding image", err)
	}
	_ = inputFile.Close()

	// generate the signed distance field image
	outputImage := sdf.GenerateDistanceFieldImage(inputImage, downscale, spread, channelFlags, threshold)

	// create output file
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		exitWithError("While creating output file", err)
	}

	// encode output image as png
	err = png.Encode(outputFile, outputImage)
	if err != nil {
		exitWithError("While encoding output image", err)
	}

	// check for any error on closing the output file
	err = outputFile.Close()
	if err != nil {
		exitWithError("While closing output file", err)
	}

	fmt.Println(" done.")
}

func pathIsFolder(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

func pathExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func exitWithError(text string, err error) {
	if err != nil {
		fmt.Println(fmt.Sprintf("\nError: %s (%s)\n\n", text, err))
	} else {
		fmt.Printf("\nError: %s\n\n", text)
	}
	os.Exit(-1)
}
