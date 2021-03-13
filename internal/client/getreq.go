package client

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type GetOpts struct {
	Query string
}

func ProcessGet(get GetOpts) error {
	fmt.Println("WE MADE IT")
	fmt.Println(serverAddress)

	resp, err := client.Get(serverAddress)
	if err != nil {
		log.Println(err)
		return err
	}

	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("%v\n", resp.Status)
	fmt.Printf(string(htmlData))
	fmt.Println(get.Query)
	return nil
}
