package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"gopkg.in/yaml.v3"
	"os/exec"
	"regexp"
	"time"
)


var revision string = "1.0"

type Listener struct {
    Containers []struct {
        Name string `yaml:"name"`
        Regexp []string `yaml:"regexp"`
        Label string `yaml:"label"`
    } `yaml:"containers"`
}

type Options struct {
	Config string  `short:"f" long:"file" env:"CONF" default:"listener.yml" description:"config file"`
}

type FondMessages struct {
    Container []ContainerMessages
}

type ContainerMessages struct {
    Name string
    Messages []string
}

func main() {
	fmt.Printf("Listener %s\n", revision)

	var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    listener, errYaml := LoadConfig(opts.Config)
    if errYaml != nil {
        log.Println(errYaml)
    }

    for _, container := range listener.Containers {
        var containerMessages ContainerMessages;
        containerMessages.Name = string(container.Name)
       // cmd := exec.Command("docker", "logs", string(container.Name), "--tail", "30")
        time := time.Now().Add(-time.Minute * 5).Format("2006-01-02T15:04:05Z")
        cmd := exec.Command("docker", "logs", string(container.Name), "--since", time)

        output, err := cmd.CombinedOutput()

        if err != nil {
            log.Fatal(err)
        }
        outStr := string(output)

        for _, regExpStr := range container.Regexp {
            matched := regexp.MustCompile(regExpStr)
            matches := matched.FindAllStringSubmatch(outStr, -1)
            matchesIndexes := matched.FindAllStringSubmatchIndex(outStr, -1)

            for _, v := range matches {
                containerMessages.Messages = append(containerMessages.Messages, v[1])
                fmt.Println(v[1])
            }
            fmt.Println(matchesIndexes)
        }

        fmt.Println(containerMessages)
        fmt.Println(container.Name)
    }
}

func LoadConfig(file string) (*Listener, error) {
	fh, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("can't load config file %s: %w", file, err)
	}
	defer fh.Close()

	res := Listener{}
	if err := yaml.NewDecoder(fh).Decode(&res); err != nil {
		return nil, fmt.Errorf("can't parse config: %w", err)
	}
	return &res, nil
}

