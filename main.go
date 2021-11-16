package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tag",
				Aliases: []string{"t"},
				Value:   "",
				Usage:   "tag of fluentd",
			},
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"h"},
				Value:   "127.0.0.1",
				Usage:   "destination host",
			},
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   "",
				Usage:   "destination port",
			},
		},
		Name:   "fluent-cat-go",
		Usage:  "fluent-cat-go -h [HOST] -p [PORT]",
		Action: readAndPost,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func readAndPost(c *cli.Context) error {
	logger, err := fluent.New(fluent.Config{})
	if err != nil {
		return xerrors.Errorf("failed to create fluent logger: %w", err)
	}
	defer logger.Close()
	tag := "myapp.access"

	reader := bufio.NewReaderSize(os.Stdin, 2<<18)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("failed to read line", err)
			continue
		}
		if isPrefix {
			log.Println("warn: msg too long")
			continue
		}
		payload := make(map[string]interface{})
		if err = json.Unmarshal(line, &payload); err != nil {
			log.Println("failed to parse json", err)
			continue
		}
		if err := logger.Post(tag, payload); err != nil {
			log.Println("failed to post", err)
			continue
		}
	}
	return nil
}
