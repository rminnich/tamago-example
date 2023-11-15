package main

import (
	"log"
	"strings"
)

func runone(s string) error {
	args := strings.Split(s, " ")
	log.Printf("run %v", args)
	return nil
}
