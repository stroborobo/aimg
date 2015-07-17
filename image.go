package aimg

import (
	"image"
	"io"
	"os"
	"strings"

	// supporting these files types
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/Knorkebrot/ansirgb"
)

var (
	isUTF8 = false
)

func init() {
	isUTF8 = strings.Contains(os.Getenv("LC_ALL"), "UTF-8") ||
		strings.Contains(os.Getenv("LANG"), "UTF-8")
}

// Image represents an ANSI color code "image"
type Image struct {
	blocks []*Block

	// actual size
	width  int
	height int

	// display size
	cols int
	rows int
}

// NewImage returns an Image that will parse and print for a given width of
// columns.
func NewImage(cols int) *Image {
	return &Image{
		cols: cols,
	}
}

// Size returns the resulting size in columns and rows.
func (im *Image) Size() (rows, cols int) {
	return im.cols, im.rows / 2
}

// ActualSize returns the size of the underlying image.
func (im *Image) ActualSize() (height, width int) {
	return im.width, im.height
}

// ParseFile is a shorthand for os.Open() and ParseReader().
func (im *Image) ParseFile(fpath string) error {
	fh, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer fh.Close()
	return im.ParseReader(fh)
}

// ParseReader reads image data from the reader and decodes it into Blocks.
func (im *Image) ParseReader(rd io.Reader) error {
	img, _, err := image.Decode(rd)
	if err != nil {
		return err
	}

	im.width = img.Bounds().Dx()
	im.height = img.Bounds().Dy()

	if im.width < im.cols {
		im.cols = im.width
	}

	ratio := float64(im.width) / float64(im.cols)
	im.rows = int(float64(im.height) / ratio)

	for r := 1; r < im.rows; r += 2 {
		for c := 1; c < im.cols; c++ {
			x := int(ratio * float64(c))
			yt := int(ratio * float64(r-1))
			yb := int(ratio * float64(r))

			b := &Block{
				Top:    ansirgb.Convert(img.At(x, yt)),
				Bottom: ansirgb.Convert(img.At(x, yb)),
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

// Blank returns a string containing as many newlines as needed to display the
// image.
func (im *Image) Blank() string {
	return strings.Repeat("\n", im.rows/2)
}

// BlankReset returns a string like Blank() but with an escape code to set the
// cursor to the first of the previous newlines.
func (im *Image) BlankReset() string {
	ret := im.Blank()
	return ret + cursorUp(len(ret))
}

// String returns the Image's string representation.
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

// Block represents two pixels or a character in a string. It contains a
// Unicode UPPER HALF BLOCK, so the top "pixel" is the foreground color and the
// bottom "pixel" is the background color.
type Block struct {
	nocolor bool
	Top     *ansirgb.Color
	Bottom  *ansirgb.Color
}

// String returns the string representation of the block. If aimg can't
// determine whether this is a UTF-8 environment, String will use a '#' instead
// of the UPPER HALF BLOCK.
func (b *Block) String() string {
	if b.nocolor {
		return " "
	}

	ret := b.Bottom.Bg()

	// Foreground colors are lighter in some terminals.
	// Ignore top (FG) if it's the same color anyway
	if b.Top.Code != b.Bottom.Code {
		ret += b.Top.Fg()
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
	return b.Bottom.Code == b2.Bottom.Code && b.Top.Code == b2.Top.Code
}
