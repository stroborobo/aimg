package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/monochromegane/terminal"
	flag "github.com/ogier/pflag"
	"github.com/stroborobo/aimg"
	aimgterm "github.com/stroborobo/aimg/terminal"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-w (num | num%% | .num) ] file [file...]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var widthstr string
	flag.StringVarP(&widthstr, "width", "w", "100%",
		"Output width. Supports column count or percentage and decimals relative to the terminal's width")
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	width, err := aimgterm.GetColumns(widthstr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "\nYou may need to set width manually using -w num")
		os.Exit(1)
	}
	width -= 1 // -1 for the reset column
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
