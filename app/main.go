package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
)


var revision string = "1.0"

type Options struct {
	Config string  `short:"f" long:"file" env:"CONF" default:"listener.yml" description:"config file"`
}

func main() {
	fmt.Printf("Listener %s\n", revision)

	var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf(opts.Config)
}
