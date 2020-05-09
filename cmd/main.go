package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"os"
	"github.com/urfave/cli/v2"
)

var(
	suffix string = ""
	comment string = ""
)

func writeEslintIgnore(p string) error{
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return err
	}

	for _, f := range files{
		if !f.IsDir() && strings.HasSuffix(f.Name(), suffix) {
			err := writeComment(p + f.Name())
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Printf("%s %s\n", "+", f.Name())
		}
	}
	return nil
}

func writeComment(p string) error {
	f, err := os.OpenFile(p, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create buffer
	buf := new(bytes.Buffer)
	buf.WriteString(comment+"\n")

	_, err = io.Copy(buf, f)
	if err != nil {
		return err
	}

	// Truncate file
	err = f.Truncate(0)
	if err != nil {
		return err
	}

	// Write content
	_,err = io.Copy(f, buf)
	if err != nil {
		return err
	}
	return nil
}


func main() {
	app := cli.NewApp()
	app.Usage = "react eslint crunch, if you use grpc-web"
	app.Description = "Add eslint ignore comment to generated files"

	app.Flags = []cli.Flag {
		&cli.StringFlag{
			Name: "suffix",
			Value: ".js",
			Aliases: []string{"s"},
			Usage: "Suffix for files name",
			Destination: &suffix,
		},
		&cli.StringFlag{
			Name: "comment",
			Value: "/* eslint-disable */",
			Aliases: []string{"c"},
			Usage: "Comment to add",
			Destination: &comment,
		},
	}

	app.Action = func(c *cli.Context) error {
		p := c.Args().Get(0)
		if p == "" {
			return cli.ShowAppHelp(c)
		}
		return writeEslintIgnore(p)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}