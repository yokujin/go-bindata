package main

import (
  "fmt"
  "io"
  "bytes"
  "os"
)

type NullCompressor struct {}

func (c *NullCompressor) Compress(in io.Reader) io.Writer {
  var out bytes.Buffer
  var written, err = io.Copy(&out, in)
  fmt.Fprintln(os.Stderr, "Written: ", written, "  Err: ", err)
  return &out
}

var in string = "Hello from Hell"

func main() {
  c := new(NullCompressor)
  fmt.Println( c.Compress(bytes.NewReader([]byte(in))) )
}
