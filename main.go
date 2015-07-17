package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	flag "github.com/ogier/pflag"
)

func main() {
	var widthstr string
	flag.StringVarP(&widthstr, "width", "w", "100%", "Output width. Supports column count and percentage.")
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	width := getColumns(widthstr) - 1 // -1 for the reset column
	for _, fpath := range flag.Args() {
		im := &Image{}
		handleErr(im.Parse(fpath, width))

		im.PrintBlank()
		im.Print()

		fmt.Println("File:", filepath.Base(fpath), "size:", im.Width, "x", im.Height)
	}
}
