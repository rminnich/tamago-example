package main

import (
	"io"
	"log"
	"strings"
)

func runone(s string, w io.Writer) error {
	args := strings.Split(s, " ")
	log.Printf("run %v", args)
	return nil
}
