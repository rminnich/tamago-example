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
	syscall.MkDev("/dev/pipe", 0666, openPipe)
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

	if hasUSB || hasEth {
		network.SetupStaticWebAssets(cmd.Banner)
		network.StartInterruptHandler(usb, eth)
	} else {
		for {
			err := doit(console)
			devcons := cmd.Console()
			term := term.NewTerminal(devcons, "uroot")
			term.SetPrompt("uroot>")

			fd, err := os.OpenFile("/dev/pipe", os.O_RDONLY, 0)
			if err != nil {
				log.Printf("pipe: %v", err)
				continue
			}
			os.Stdin = fd
			w, err := os.OpenFile("/dev/pipe", os.O_WRONLY, 0)
			if err != nil {
				log.Printf("write pipe: %v", err)
				continue
			}
			runtime.Exit = exited
			s, err := term.ReadLine()
			log.Printf("readline %q %v", s, err)

			if err == io.EOF {
				continue
			}
			if err != nil {
				log.Printf("readline error, %v", err)
				continue
			}

			term.SetPrompt("")
			go func() {
				for {
					s, err := term.ReadLine()
					log.Printf("readline %q %v", s, err)

					if err == io.EOF || len(s) == 0 {
						log.Printf("%v %v EOF", err, len(s))
						if _, err := w.Write([]byte{}); err != nil {
							log.Printf("pipe write:%v", err)
						}
						return
					}

					if err != nil {
						log.Printf("readline error, %v", err)
						if _, err := w.Write([]byte{}); err != nil {
							log.Printf("pipe write:%v", err)
						}
						continue
					}

					if _, err := w.Write([]byte(s + "\n")); err != nil {
						log.Printf("pipe write:%v", err)
					}
				}
			}()
			log.Printf("run %s", s)
			runone(s, w)
			log.Printf("runoene done")
		}

	}
}

func exited() {
	panic("exiting")
}

func doit(console *cmd.Interface) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered. Error:%v\n", r)
			err = fmt.Errorf("wel: %v", r)
		}
	}()
	cmd.SerialConsole(console)
	return err
}
