package main

import (
  "bytes"
)

type Compressor interface {
  Compress(out io.Writer, in io.Reader)
}

