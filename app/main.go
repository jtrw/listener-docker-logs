package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"gopkg.in/yaml.v3"
	"os/exec"
	"regexp"
)


var revision string = "1.0"

type Listener struct {
    Containers []struct {
        Name string `yaml:"name"`
        Regexp string `yaml:"regexp"`
        Label string `yaml:"label"`
    } `yaml:"containers"`
}

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

    //fmt.Printf(opts.Config)

    listener, errYaml := LoadConfig(opts.Config)
    if errYaml != nil {
        log.Println(errYaml)
    }

    for _, container := range listener.Containers {
        cmd := exec.Command("docker", "logs", string(container.Name), "--tail", "10")

        output, err := cmd.CombinedOutput()

        if err != nil {
            log.Fatal(err)
        }
        outStr := string(output)

        matched := regexp.MustCompile(container.Regexp)
        fmt.Println(matched.FindStringIndex(string(outStr)))

        //fmt.Println(string(output))

        //fmt.Println(string(out))

        fmt.Println(container.Name)
    }
    //fmt.Println(listener)
}

func LoadConfig(file string) (*Listener, error) {
	fh, err := os.Open(file) //nolint
	if err != nil {
		return nil, fmt.Errorf("can't load config file %s: %w", file, err)
	}
	defer fh.Close() //nolint

	res := Listener{}
	if err := yaml.NewDecoder(fh).Decode(&res); err != nil {
		return nil, fmt.Errorf("can't parse config: %w", err)
	}
	return &res, nil
}

