package bbtamago

import bbmain "bb.u-root.com/bb/pkg/bbmain"
import "log"
import "os"
import "flag"
import "strings"

func runone(s string) error {
	os.Args = strings.Split(s, " ")
	log.Printf("run %v", os.Args)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	return bbmain.Run(os.Args[0])
}
