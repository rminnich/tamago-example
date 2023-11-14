// https://github.com/usbarmory/tamago-example
//
// Copyright (c) WithSecure Corporation
// https://foundry.withsecure.com
//
// Use of this source code is governed by the license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/usbarmory/tamago/soc/nxp/enet"
	"github.com/usbarmory/tamago/soc/nxp/usb"
	"golang.org/x/term"

	"github.com/usbarmory/tamago-example/cmd"
	"github.com/usbarmory/tamago-example/network"
)

var Build string
var Revision string

// can not initialize this. It gets turned into a function
// that runs after init()
var onexit func()

func init() {
	log.SetFlags(0)

	cmd.Banner = fmt.Sprintf("%s/%s (%s) • %s %s",
		runtime.GOOS, runtime.GOARCH, runtime.Version(),
		Revision, Build)

	cmd.Banner += fmt.Sprintf(" • %s", cmd.Target())
}

func main() {
	var usb *usb.USB
	var eth *enet.ENET

	logFile, _ := os.OpenFile("/tamago-example.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	hasUSB, hasEth := cmd.HasNetwork()

	console := &cmd.Interface{}

	if hasUSB {
		usb = network.StartUSB(console.Start, logFile)
	}

	if hasEth {
		eth = network.StartEth(console.Start, logFile)
	}

	cmd.NIC = eth

	var fd [2]int
	syscall.Pipe(fd[:])
	if err := syscall.Dup2(0, fd[0]); err != nil {
		log.Printf("syscall.Dup2(0, %d): %v", fd[0], err)
	}
	os.Stdin = os.NewFile(0, "pipe input to os.Stdin")
	w := os.NewFile(uintptr(fd[1]), "pipe output to os.Stdin")

	if hasUSB || hasEth {
		network.SetupStaticWebAssets(cmd.Banner)
		network.StartInterruptHandler(usb, eth)
	} else {
		for {
			err := doit(console)
			log.Printf("console err %v; take 5", err)
			go func() {
				devcons := cmd.Console()
				term := term.NewTerminal(devcons, "uroot")
				term.SetPrompt("uroot>")

				for {
					s, err := term.ReadLine()
					log.Printf("readline %q %v", s, err)

					if err == io.EOF {
						return
					}

					if err != nil {
						log.Printf("readline error, %v", err)
						continue
					}

					if _, err := w.Write([]byte(s)); err != nil {
						log.Printf("pipe write:%v", err)
					}
				}
			}()
			time.Sleep(5 * time.Second)
		}
	}
}

func doit(console *cmd.Interface) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered. Error:%v\n", r)
			err = fmt.Errorf("wel: %v", r)
		}
	}()
	runtime.Exit = func() {
		panic("hey we exited")
	}
	cmd.SerialConsole(console)
	return err
}
