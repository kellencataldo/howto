package client

import (
	"fmt"
)

type GetOpts struct {
	Query string
}

func ProcessGet(opts HowtoOpts, get GetOpts) error {
	fmt.Println(get.Query)
	return nil
}
