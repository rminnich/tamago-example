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

var pipe *pipeFile

// This is a real kludge for now but we'll do more later.

func openPipe() (syscall.DevFile, error) {
	if pipe == nil {
		r, w := io.Pipe()
		pipe = &pipeFile{r: r, w: w}
	}
	return pipe, nil
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
