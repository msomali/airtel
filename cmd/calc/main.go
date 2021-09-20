package main

import (
	"github.com/techcraftt/base/cli"
	"log"
	"os"
)

func main() {
	client := cli.New("localhost", 9099)
	if err := client.Run(os.Args); err != nil {
		log.Fatalf("err: %v\n", err)
	}
}
