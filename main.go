package main

import (
	"log"

	"github.com/revolutionchain/charon/cli"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	cli.Run()
}
