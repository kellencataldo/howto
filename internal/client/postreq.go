package client

import (
	"fmt"
	"os"
)

type PostOpts struct {
	File  os.File
	Title string
}

func ProcessPost(opts HowtoOpts, post PostOpts) error {
	fmt.Println(post.Title)
	return nil
}
