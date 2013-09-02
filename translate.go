// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package main

import (
	"compress/gzip"
	"fmt"
	"io"
)

// translate translates the input file to go source code.
func translate(input io.Reader, output io.Writer, pkgname, funcname string, uncompressed) {
  if uncompressed {
    translate_memcpy_uncomp(input, output, pkgname, funcname)
  } else {
    translate_memcpy_comp(input, output, pkgname, funcname)
  }
}

// input -> gzip -> gowriter -> output.
func translate_memcpy_comp(input io.Reader, output io.Writer, pkgname, funcname string) {
	fmt.Fprintf(output, `package %s

import (
	"bytes"
	"compress/gzip"
	"io"
)

// %s returns raw, uncompressed file data.
func %s() []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte{`, pkgname, funcname, funcname)

	gz := gzip.NewWriter(&ByteWriter{Writer: output})
	io.Copy(gz, input)
	gz.Close()

	fmt.Fprint(output, `
	}))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes()
}`)
}

// input -> gzip -> gowriter -> output.
func translate_memcpy_uncomp(input io.Reader, output io.Writer, pkgname, funcname string) {
	fmt.Fprintf(output, `package %s

// %s returns raw file data.
func %s() []byte {
	return []byte{`, pkgname, funcname, funcname)

	io.Copy(&ByteWriter{Writer: output}, input)

	fmt.Fprint(output, `
	}
}`)
}

