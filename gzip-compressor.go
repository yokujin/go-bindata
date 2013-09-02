package main

import (
  "fmt"
  "io"
  "bytes"
  "os"
  "compress/gzip"
)

type GzipCompressor struct {}

func (c *GzipCompressor) Compress(in io.Reader) io.Writer {
  var out bytes.Buffer
  var gz = gzip.NewWriter(&out)
  var written, err = io.Copy(gz, in)
  gz.Close()
  fmt.Fprintln(os.Stderr, "Written: ", written, "  Err: ", err)
  return &out
}

var in string = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func main() {
  c := new(GzipCompressor)
  fmt.Print( c.Compress(bytes.NewReader([]byte(in))) )
}
