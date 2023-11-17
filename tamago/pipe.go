// https://github.com/usbarmory/tamago-example
//
// Copyright (c) WithSecure Corporation
// https://foundry.withsecure.com
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

package main

import (
	"io"
	"syscall"
)

type pipeFile struct {
	r *io.PipeReader
	w *io.PipeWriter
}

func openPipe() (syscall.DevFile, error) {
	r, w := io.Pipe()
	return pipeFile{r: r, w: w}, nil
}

func (f pipeFile) close() error {
	f.w.Close()
	return nil
}

func (f pipeFile) Pread(b []byte, offset int64) (int, error) {
	n, err := f.r.Read(b)
	return n, err
}

func (f pipeFile) Pwrite(b []byte, offset int64) (int, error) {
	n, err := f.w.Write(b)
	return n, err
}
