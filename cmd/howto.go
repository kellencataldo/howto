package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/kellencataldo/howto/internal/client"
	log "github.com/sirupsen/logrus"
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

func PopulateMissingConnFields(conn *client.ConnectionInfo) error {

	// @TODO: check environment variables here, if one is not found return error
	/*
		type ConnectionInfo struct {
			ServerAddress string
			ServerPort int
			ClientCert os.File
			ClientKey  os.File
			ServerCert os.File
	*/

	return nil
}

func main() {

	log.SetOutput(os.Stdout)
	// @TODO make this customizable
	// ie: var mySelector *string = parser.Selector("d", "debug-level", []string{"INFO", "DEBUG", "WARN"}, ...)
	log.SetLevel(log.WarnLevel)

	if len(os.Args) > 1 && !flagsPresent(os.Args[1:]) {

		conn := client.ConnectionInfo{}
		if err := PopulateMissingConnFields(&conn); err != nil {
			log.Fatalln(err)
		} else if err = client.InitializeConnectionInfo(conn); err != nil {
			log.Fatalln(err)
		}

		client.ProcessGet(client.GetOpts{Query: strings.Join(os.Args[1:], " ")})
	}

	parser := argparse.NewParser("howto", "takes query string and returns a how-to article")
	address := parser.String("a", "address", &argparse.Options{Required: false, Help: "address of the howto server, will also search environment variable HOWTO_ADDR"})
	port := parser.Int("p", "port", &argparse.Options{Required: false, Help: "port used to connect to server, will also seach environment variable HOWTO_PORT"})
	cert := parser.File("c", "cert", os.O_RDONLY, 0400, &argparse.Options{Required: false, Help: "client certificate, will also search environment variable HOWTO_CERT"})
	key := parser.File("k", "key", os.O_RDONLY, 0400, &argparse.Options{Required: false, Help: "client key, will also search environment variable HOWTO_KEY"})
	srvCert := parser.File("s", "srvcert", os.O_RDONLY, 0400, &argparse.Options{Required: false, Help: "server cert, will also search environment variable HOWTO_SRV_CERT"})

	getCmd := parser.NewCommand("get", "gets articles from the server")
	query := getCmd.String("q", "query", &argparse.Options{Required: true, Help: "search string to use"})

	postCmd := parser.NewCommand("post", "posts an article to the server")
	file := postCmd.File("f", "file", os.O_RDONLY, 0400, &argparse.Options{Required: true, Help: "file to post"})
	title := postCmd.String("t", "title", &argparse.Options{Required: true, Help: "title of the article, will be matched against during queries"})

	err := parser.Parse(os.Args)
	if nil != err {
		log.Fatalln(parser.Usage(err))
	}

	/*
		type ConnectionInfo struct {
			ServerAddress string
			ServerPort int
			ClientCert os.File
			ClientKey  os.File
			ServerCert os.File
	*/

	conn := client.ConnectionInfo{}
	if *address != "" {
		conn.ServerAddress = *address
	}

	if *port != 0 {
		conn.ServerPort = *port
	}

	if cert != nil {
		conn.ClientCert = cert
	}

	if key != nil {
		conn.ClientKey = key
	}

	if srvCert != nil {
		conn.ServerCert = srvCert
	}

	if err := PopulateMissingConnFields(&conn); err != nil {
		log.Fatalln(err)
	} else if err = client.InitializeConnectionInfo(conn); err != nil {
		log.Fatalln(err)
	}

	// kellen, initialize global variables here.
	if getCmd.Happened() {
		get := client.GetOpts{Query: *query}
		err = client.ProcessGet(get)
	} else if postCmd.Happened() {
		post := client.PostOpts{File: *file, Title: *title}
		err = client.ProcessPost(post)
	} else {
		log.Fatalln(parser.Usage(fmt.Sprintf("Unable to process request: %s", strings.Join(os.Args, " "))))
	}

	if nil != err {
		log.Fatalln(err)
	}

	os.Exit(0)
}
