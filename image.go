package aimg

import (
	"bytes"
	"image"
	"io"
	"os"
	"strings"

	// supporting these files types
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/stroborobo/aimg/terminal"
	"github.com/stroborobo/ansirgb"
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
	im image.Image

	// actual size
	width  int
	height int

	// display size
	cols int
	rows int

	ratio float64

	readIndex int
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

// ParseReader reads image data from the reader.
func (im *Image) ParseReader(rd io.Reader) error {
	if img, _, err := image.Decode(rd); err != nil {
		return err
	} else {
		im.im = img
	}

	im.width = im.im.Bounds().Dx()
	im.height = im.im.Bounds().Dy()

	if im.width < im.cols {
		im.cols = im.width
	}

	im.ratio = float64(im.width) / float64(im.cols)
	im.rows = int(float64(im.height) / im.ratio)

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
	return ret + terminal.CursorUp(len(ret))
}

// WriteTo writes the image data to wr.
func (im *Image) WriteTo(wr io.Writer) (int, error) {
	written := 0
	for r := 0; r < im.rows; r += 2 {
		var before *Block
		for c := 0; c < im.cols; c++ {
			x := int(im.ratio * float64(c))
			yt := int(im.ratio * float64(r))
			yb := int(im.ratio * float64(r+1))

			b := &Block{
				Top:    im.getColor(x, yt),
				Bottom: im.getColor(x, yb),
			}

			if before != nil {
				b.nocolor = b.equals(before)
			}
			before = b

			n, err := io.WriteString(wr, b.String())
			written += n
			if err != nil {
				return written, err
			}
		}
		n, err := io.WriteString(wr, terminal.NewLine())
		written += n
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

func (im *Image) getColor(x, y int) *ansirgb.Color {
	c := im.im.At(x, y)
	if _, _, _, a := c.RGBA(); a < 1<<14-1 {
		return &ansirgb.Color{c, -1} // default/transparent color
	}
	return ansirgb.Convert(c)
}

// String returns the Image's string representation. It's a shorthand to
// WriteTo() using a bytes.Buffer and buf.String().
func (im *Image) String() string {
	buf := &bytes.Buffer{}
	im.WriteTo(buf)
	return buf.String()
}

// Block represents two pixels or a character in a string. It contains a
// Unicode LOWER/UPPER HALF BLOCK, so the top or bottom "pixel" is the
// foreground color and the other one is the background color.
// The character used might change if transparency is needed, which is only
// possible by using the default background color.
type Block struct {
	nocolor bool
	Top     *ansirgb.Color
	Bottom  *ansirgb.Color
}

// String returns the string representation of the block. If aimg can't
// determine whether this is a UTF-8 environment, String will use a '#' instead
// of the LOWER/UPPER HALF BLOCK. If the block's color is equal to the one
// before, it'll just return a string containing a single space.
func (b *Block) String() string {
	// By default top is background, bottom foreground, using
	// LOWER HALF BLOCK.
	first := b.Top
	second := b.Bottom
	block := "\u2584"

	// If the bottom is transparent however switch them all
	// since transparency can only be archieved by using the default
	// background color.
	if b.Bottom.IsTransparent() {
		first, second = second, first
		block = "\u2580"
	}

	ret := ""
	if !b.nocolor {
		ret += first.Bg()
	}

	// Foreground colors are lighter in some terminals. Also "transparent"
	// foreground is not a thing.  Ignore it if it's the same color anyway.
	if !first.Equals(second) {
		if !b.nocolor {
			ret += second.Fg()
		}
		// If it's not a UTF-8 terminal, fall back to '#'
		if isUTF8 {
			ret += block
		} else {
			ret += "#"
		}
	} else {
		ret += " "
	}

	return ret
}

func (b *Block) equals(b2 *Block) bool {
	return b.Bottom.Equals(b2.Bottom) && b.Top.Equals(b2.Top)
}
