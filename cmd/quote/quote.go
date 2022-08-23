package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func usage() string {
	return fmt.Sprintf(`Outputs input as a quoted string.

Usage:
  %s [flags] [file]

Flags:
	-h, --help		Displays this message
	-v, --verbose		Be verbose

`, os.Args[0])

}

func exitf(code int, format string, a ...interface{}) {
	w := os.Stderr
	if code == 0 {
		w = os.Stdout
	}

	fmt.Fprintf(w, format, a...)
	os.Exit(code)
}

type options struct {
	verbose bool
	in      io.Reader
}

func mustParseOpts() *options {
	opts := &options{
		in: os.Stdin,
	}

	var fileopened bool

	var args []string
	for _, arg := range os.Args[1:] { // convert '--foo=blah' into '--foo blah' for parsing
		args = append(args, strings.SplitN(arg, "=", 2)...)
	}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h", "--help":
			exitf(0, usage())
		case "-v", "--verbose":
			opts.verbose = true
		default:
			if fileopened {
				// opts.in.Close()
				exitf(-1, "multiple files supplied\n")
			}
			fileopened = true
			f, err := os.Open(args[i])
			if err != nil {
				exitf(-1, "error opening file: %s\n", err.Error())
			}
			opts.in = f
		}
	}
	return opts
}

func main() {
	o := mustParseOpts()
	b, err := ioutil.ReadAll(o.in)
	if err != nil {
		exitf(-1, "error reading input: %s\n", err.Error())
	}
	fmt.Printf("%q\n", string(b))
}
