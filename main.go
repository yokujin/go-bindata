// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var (
	pipe         = false
	in           = ""
	out          = flag.String("out", "", "Optional path and name of the output file.")
	pkgname      = flag.String("pkg", "main", "Name of the package to generate.")
	funcname     = flag.String("func", "", "Optional name of the function to generate.")
	prefix       = flag.String("prefix", "", "Optional path prefix to strip off map keys and function names.")
	uncompressed = flag.Bool("uncompressed", false, "The specified resource will /not/ be GZIP compressed when this flag is specified. This alters the generated output code.")
	toc          = flag.Bool("toc", false, "Generate a table of contents for this and other files. The input filepath becomes the map key. This option is only useable in non-pipe mode.")
	version      = flag.Bool("version", false, "Display version information.")
	regFuncName  = regexp.MustCompile(`[^a-zA-Z0-9_]`)
)

func main() {
	parseArgs()

	if pipe {
		translate(os.Stdin, os.Stdout, *pkgname, *funcname, *uncompressed)
		return
	}

	fs, err := os.Open(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[e] %s\n", err)
		return
	}

	defer fs.Close()

	fd, err := os.Create(*out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[e] %s\n", err)
		return
	}

	defer fd.Close()

	// Translate binary to Go code.
	translate(fs, fd, *pkgname, *funcname, *uncompressed)

	// Append the TOC init function to the end of the output file and
	// write the `bindata-toc.go` file, if applicable.
	if *toc {
		dir, _ := filepath.Split(*out)
		err := createTOC(dir, *pkgname)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s\n", err)
			return
		}

		writeTOCInit(fd, in, *prefix, *funcname)
	}
}

// parseArgs processes and verifies commandline arguments.
func parseArgs() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <filename>\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	pipe = flag.NArg() == 0

	if !pipe {
		*prefix, _ = filepath.Abs(filepath.Clean(*prefix))
		in, _ = filepath.Abs(filepath.Clean(flag.Args()[0]))
		*out = safeFilename(*out, in)
	}

	if len(*pkgname) == 0 {
		fmt.Fprintln(os.Stderr, "[w] No package name specified. Using 'main'.")
		*pkgname = "main"
	} else {
		if unicode.IsDigit(rune((*pkgname)[0])) {
			// Identifier can't start with a digit.
			*pkgname = "_" + *pkgname
		}
	}

	if len(*funcname) == 0 {
		if pipe {
			// Can't infer from input file name in this mode.
			fmt.Fprintln(os.Stderr, "[e] No function name specified.")
			os.Exit(1)
		}

		*funcname = safeFuncname(in, *prefix)
		fmt.Fprintf(os.Stderr, "[w] No function name specified. Using %s.\n", *funcname)
	}
}

// safeFilename creates a safe output filename from the given
// output and input paths.
func safeFilename(out, in string) string {
	var filename string

	if len(out) == 0 {
		filename = in + ".go"

		_, err := os.Lstat(filename)
		if err == nil {
			// File already exists. Pad name with a sequential number until we
			// find a name that is available.
			count := 0

			for {
				filename = path.Join(out, fmt.Sprintf("%s.%d.go", in, count))
				_, err = os.Lstat(filename)

				if err != nil {
					break
				}

				count++
			}
		}
	} else {
		filename, _ = filepath.Abs(filepath.Clean(out))
	}

	// Ensure output directory exists while we're here.
	dir, _ := filepath.Split(filename)
	_, err := os.Lstat(dir)
	if err != nil {
		os.MkdirAll(dir, 0755)
	}

	return filename
}

// safeFuncname creates a safe function name from the input path.
func safeFuncname(in, prefix string) string {
	name := strings.Replace(in, prefix, "", 1)

	if len(name) == 0 {
		name = in
	}

	name = strings.ToLower(name)
	name = regFuncName.ReplaceAllString(name, "_")

	if unicode.IsDigit(rune(name[0])) {
		// Identifier can't start with a digit.
		name = "_" + name
	}

	// Get rid of "__" instances for niceness.
	for strings.Index(name, "__") > -1 {
		name = strings.Replace(name, "__", "_", -1)
	}

	// Leading underscore is silly.
	if name[0] == '_' {
		name = name[1:]
	}

	return name
}
