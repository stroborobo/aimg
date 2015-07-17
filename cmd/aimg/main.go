package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Knorkebrot/aimg"
	"github.com/monochromegane/terminal"
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
		im := aimg.NewImage(width)
		handleErr(im.ParseFile(fpath))

		if terminal.IsTerminal(os.Stdout) {
			fmt.Print(im.BlankReset())
		}
		fmt.Print(im)

		w, h := im.ActualSize()
		fmt.Println("File:", filepath.Base(fpath), "size:", w, "x", h)
	}
}
