package main

import (
	"fmt"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"github.com/Knorkebrot/ansirgb"
	"github.com/olekukonko/ts"
)

type Block struct {
	top	*ansirgb.Color
	bottom	*ansirgb.Color
}

func (b *Block) String() string {
	ret := fmt.Sprintf("\033[48;5;%dm", b.bottom.Code)
	if b.top != nil {
		ret += fmt.Sprintf("\033[38;5;%dm", b.top.Code)
		// If it's not a UTF-8 terminal, fall back to '#'
		if strings.Contains(os.Getenv("LC_ALL"), "UTF-8") ||
		   strings.Contains(os.Getenv("LANG"), "UTF-8") {
			ret += "\u2580"
		} else {
			ret += "#"
		}
	} else {
		ret += " "
	}
	return ret
}

func reset() {
	// add a space to prevent artifacts after resizing
	fmt.Printf("\033[0m ")
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s file [file...]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var width int
	flag.IntVar(&width, "w", 0, "Output width, use 0 for terminal width")
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if width == 0 {
		size, err := ts.GetSize()
		if err != nil {
			fmt.Fprintln(os.Stderr, err, "\nYou may need to "+
				"set width manually using -w num")
			os.Exit(2)
		}
		width = size.Col() - 1	// -1 for the reset column
	}

	for _, fpath := range flag.Args() {
		fh, err := os.Open(fpath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(10)
		}

		img, _, err := image.Decode(fh)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			fh.Close()
			os.Exit(20)
		}

		imgWidth := img.Bounds().Dx()
		imgHeight := img.Bounds().Dy()

		if imgWidth < width {
			width = imgWidth
		}

		ratio := float64(imgWidth) / float64(width)
		rows := int(float64(imgHeight) / ratio)

		for i := 1; i < rows; i += 2 {
			for j := 0; j < width; j++ {
				// TODO: get average color of the area instead
				// of one pixel?
				x := int(ratio * float64(j))
				yTop := int(ratio * float64(i - 1))
				yBottom := int(ratio * float64(i))

				top := ansirgb.Convert(img.At(x, yTop))
				bottom := ansirgb.Convert(img.At(x, yBottom))

				b := &Block{}
				b.bottom = bottom

				// Foreground colors are lighter in some terminals.
				// Ignore top (FG) if it's the same color anyway
				if top.Code != bottom.Code {
					b.top = top
				}

				fmt.Printf("%s", b)
			}
			reset()
			fmt.Printf("\n")
		}
		fh.Close()

		fmt.Println("File:", filepath.Base(fpath), "size:", imgWidth, "x", imgHeight)
	}
}
