package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/kellencataldo/howto/internal/client"
	"os"
	"strings"
	"unicode/utf8"
)

func flagsPresent(args []string) bool {

	for _, arg := range args {
		r, _ := utf8.DecodeRuneInString(arg)
		if r == '-' {
			return true
		}
	}

	return false
}

func main() {

	if len(os.Args) > 1 && !flagsPresent(os.Args[1:]) {
		htOpts := client.HowtoOpts{Address: ""}
		get := client.GetOpts{Query: strings.Join(os.Args[1:], " ")}
		client.ProcessGet(htOpts, get)
	}

	parser := argparse.NewParser("howto", "takes query string and returns a how-to article")
	a := parser.String("a", "address", &argparse.Options{Required: false, Help: "address of the howto server, will also search environment variable HOWTO_ADDR"})

	getCmd := parser.NewCommand("get", "gets articles from the server")
	q := getCmd.String("q", "query", &argparse.Options{Required: true, Help: "search string to use"})

	postCmd := parser.NewCommand("post", "posts an article to the server")
	f := postCmd.File("f", "file", os.O_RDONLY, 0400, &argparse.Options{Required: true, Help: "file to post"})
	t := postCmd.String("t", "title", &argparse.Options{Required: true, Help: "title of the article, will be matched against during queries"})

	err := parser.Parse(os.Args)
	if nil != err {
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}

	htOpts := client.HowtoOpts{Address: *a}
	if getCmd.Happened() {
		get := client.GetOpts{Query: *q}
		err = client.ProcessGet(htOpts, get)
	} else if postCmd.Happened() {
		post := client.PostOpts{File: *f, Title: *t}
		err = client.ProcessPost(htOpts, post)
	} else {
		fmt.Print(parser.Usage(fmt.Sprintf("Unable to process request: %s", strings.Join(os.Args, " "))))
		os.Exit(1)
	}

	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
