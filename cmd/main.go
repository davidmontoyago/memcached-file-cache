package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/davidmontoyago/interview-davidmontoyago-d660952eff664d8bac96c9124d7f8582/pkg/filecache"

	"github.com/urfave/cli/v2"
)

const memcachedAddr = "localhost:11211"

func main() {
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		{
			Name:  "put",
			Usage: "put a file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "filepath",
					Aliases:  []string{"f"},
					Required: true,
				},
			},
			Action: upload,
		},
		{
			Name:  "get",
			Usage: "get a file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "key",
					Aliases:  []string{"k"},
					Required: true,
				},
			},
			Action: download,
		},
	}
	app.Run(os.Args)
}

func upload(c *cli.Context) error {
	filepath := c.String("filepath")
	fmt.Printf("uploading file %s...\n", filepath)
	file, err := readFile(filepath)
	if err != nil {
		return err
	}

	mc := memcache.New(memcachedAddr)
	fc := filecache.New(mc)

	key, err := fc.Put(file)
	if err != nil {
		return err
	}

	fmt.Printf("success! file key is %s\n", key)
	return nil
}

func download(c *cli.Context) error {
	fileKey := c.String("key")
	fmt.Printf("downloading file %s...\n", fileKey)

	mc := memcache.New(memcachedAddr)
	fc := filecache.New(mc)

	file, err := fc.Get(fileKey)
	if err != nil {
		return err
	}

	localFileName := fmt.Sprintf("./%s.dat", fileKey)
	err = ioutil.WriteFile(localFileName, file, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("success! file saved to %s", localFileName)
	return nil
}

func readFile(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return file, nil
}
