package main

import (
	"fmt"
	"os"
)

func handleErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
