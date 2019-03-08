package main

import (
	"log"
	"os"
)

var (
	stdLogger   = log.New(os.Stdout, "", 0)
	errorLogger = log.New(os.Stderr, "", 0)
)

func onError(err error) {
	errorLogger.Fatal(err)
}

func print(msg string) {
	stdLogger.Println(msg)
}