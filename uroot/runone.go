package bbtamago

import bbmain "bb.u-root.com/bb/pkg/bbmain"
"log"
import "os"
import "flag"

func runone(c *Command) error {
	// put a recover here for the panic at some point.
   os.Args = append([]string{c.cmd}, c.argv...)
   flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)
   log.Printf("run %v", c)
   return bbmain.Run(c.cmd)
}
