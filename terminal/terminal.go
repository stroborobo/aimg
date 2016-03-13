package terminal

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/ts"
)

var (
	ErrNoTermSize = errors.New("Can't get terminal size.")
)

// CursorUp returns the escape code to set the cursor n lines before the
// current position.
func CursorUp(count int) string {
	return fmt.Sprintf("\033[%dA", count)
}

// NewLine returns a string containing an escape code to reset all colors and a
// space. This is needed to prevent artifacts after the end of this line if the
// terminal gets resized.
func NewLine() string {
	return "\033[0m \n"
}

// GetColumns returns the column count based on a string passed. Supported formats:
// 20% (relative to term width)
// 0.2 (relative to term width)
// 20  (columns)
func GetColumns(widthstr string) (int, error) {
	numstr := ""
	format := ""
	if strings.HasSuffix(widthstr, "%") {
		if len(widthstr) < 2 {
			fmt.Fprintf(os.Stderr, "Invalid percentage.\n")
			os.Exit(1)
		}
		numstr = widthstr[:len(widthstr)-1]
		format = "percent"
	} else if strings.Contains(widthstr, ".") {
		numstr = widthstr
		format = "decimal"
	} else {
		numstr = widthstr
		format = "columns"
	}

	num := 0
	if format == "decimal" {
		f, err := strconv.ParseFloat(numstr, 64)
		if err != nil {
			return 0, err
		}
		num = int(100.0 * f)
	} else {
		var err error
		num, err = strconv.Atoi(numstr)
		if err != nil {
			return 0, err
		}
	}

	if format == "columns" && num > 0 {
		return num, nil
	}

	size, err := ts.GetSize()
	if err != nil {
		return 0, ErrNoTermSize
	}

	if format == "columns" {
		return size.Col(), nil
	}

	cols := float64(size.Col())
	return int(cols * (float64(num) / 100.0)), nil
}
