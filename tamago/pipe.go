// https://github.com/usbarmory/tamago-example
//
// Copyright (c) WithSecure Corporation
// https://foundry.withsecure.com
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

package main

import "syscall"

type pipeFile struct{}

func openPipe() (syscall.DevFile, error) {
	return pipeFile{}, nil
}

func (f pipeFile) close() error {
	return nil
}

func (f pipeFile) Pread(b []byte, offset int64) (int, error) {
	return len(b), nil
}

func (f pipeFile) Pwrite(b []byte, offset int64) (int, error) {
	return len(b), nil
}
