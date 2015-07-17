aimg
====

[![GoDoc](https://godoc.org/github.com/Knorkebrot/aimg?status.svg)][gd]

A package to convert images to ANSI color coded text. `aimg` uses
unicode UPPER HALF BLOCK to archieve twice the pixels (tm) per
character.


Usage
-----

As a package: see [GoDoc][gd]

As a command:

	> aimg -h
	Usage: aimg file [file...]
	  -w=0: Output width, use 0 for terminal width

![screenshot](screenshot.jpg)


Get the cli tool
----------------

	> go get github.com/Knorkebrot/aimg/cmd/aimg

That's it :)

- - - -

Inspired by [minikomi's ansipix][ap].

[ap]: https://github.com/minikomi/ansipix
[gd]: https://godoc.org/github.com/Knorkebrot/aimg
