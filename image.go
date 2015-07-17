package aimg

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"strings"

	"github.com/Knorkebrot/ansirgb"
)

var (
	isUTF8 = false
)

func init() {
	isUTF8 = strings.Contains(os.Getenv("LC_ALL"), "UTF-8") ||
		strings.Contains(os.Getenv("LANG"), "UTF-8")
}

type Image struct {
	blocks []*Block
	Width  int // actual size
	Height int
	cols   int // display size
	rows   int
}

func NewImage(cols int) *Image {
	return &Image{
		cols: cols,
	}
}

func (im *Image) ParseFile(fpath string) error {
	fh, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer fh.Close()
	return im.ParseReader(fh)
}

func (im *Image) ParseReader(rd io.Reader) error {
	img, _, err := image.Decode(rd)
	if err != nil {
		return err
	}

	im.Width = img.Bounds().Dx()
	im.Height = img.Bounds().Dy()

	if im.Width < im.cols {
		im.cols = im.Width
	}

	ratio := float64(im.Width) / float64(im.cols)
	im.rows = int(float64(im.Height) / ratio)

	for r := 1; r < im.rows; r += 2 {
		for c := 1; c < im.cols; c++ {
			x := int(ratio * float64(c))
			yt := int(ratio * float64(r-1))
			yb := int(ratio * float64(r))

			b := &Block{
				top:    ansirgb.Convert(img.At(x, yt)),
				bottom: ansirgb.Convert(img.At(x, yb)),
			}

			if c > 1 {
				before := im.blocks[len(im.blocks)-1]
				b.nocolor = b.equals(before)
			}

			im.blocks = append(im.blocks, b)
		}
	}
	return nil
}

func (im *Image) Blank() string {
	c := im.rows / 2
	return strings.Repeat("\n", c)
}

func (im *Image) BlankReset() string {
	c := im.rows / 2
	return strings.Repeat("\n", c) + cursorUp(c)
}

func (im *Image) String() string {
	ret := ""
	for i, b := range im.blocks {
		if i > 0 && i%(im.cols-1) == 0 {
			ret += newLine()
		}
		ret += b.String()
	}
	return ret + newLine()
}

type Block struct {
	nocolor bool
	top     *ansirgb.Color
	bottom  *ansirgb.Color
}

func (b *Block) String() string {
	if b.nocolor {
		return " "
	}

	ret := b.bottom.Bg()

	// Foreground colors are lighter in some terminals.
	// Ignore top (FG) if it's the same color anyway
	if b.top.Code != b.bottom.Code {
		ret += b.top.Fg()
		// If it's not a UTF-8 terminal, fall back to '#'
		if isUTF8 {
			ret += "\u2580"
		} else {
			ret += "#"
		}
	} else {
		ret += " "
	}

	return ret
}

func (b *Block) equals(b2 *Block) bool {
	return b.bottom.Code == b2.bottom.Code && b.top.Code == b2.top.Code
}
