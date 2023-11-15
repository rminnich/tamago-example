package bbtamago

import bbmain "bb.u-root.com/bb/pkg/bbmain"
import "fmt"
import "log"
import "os"
import "flag"
import "strings"

func runone(s string) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			log.Printf("uroot Recovered. Error:%v\n", r)
			err = fmt.Errorf("wel: %v", r)
		}
	}()
	os.Args = strings.Split(s, " ")
	log.Printf("run %v", os.Args)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	err = bbmain.Run(os.Args[0])
	return err
}
