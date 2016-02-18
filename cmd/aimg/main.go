package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/monochromegane/terminal"
	flag "github.com/ogier/pflag"
	"github.com/stroborobo/aimg"
)

func main() {
	var widthstr string
	flag.StringVarP(&widthstr, "width", "w", "100%",
		"Output width. Supports column count, percentage and decimals.")
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
		im.WriteTo(os.Stdout)

		w, h := im.ActualSize()
		fmt.Println("File:", filepath.Base(fpath), "size:", w, "x", h)
	}
}
