package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/google/go-github/github"
	git "gopkg.in/src-d/go-git.v4"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	username := flag.String("u", "", "username")
	flag.Parse()

	if *username == "" {
		log.Fatal("Please pass in a -u flag")
	}

	ctx := context.Background()
	client := github.NewClient(nil)

	repos, _, err := client.Repositories.List(ctx, *username, nil)
	must(err)

	err = os.MkdirAll(*username, os.ModePerm)
	must(err)

	for _, repo := range repos {
		path := path.Join(*username, *repo.Name)
		url := fmt.Sprintf("https://github.com/%s/%s", *username, *repo.Name)

		_, err = git.PlainClone(path, false, &git.CloneOptions{
			URL: url,
		})
		if err != nil {
			fmt.Printf("%s: %s\n", url, err)
		}
	}
}
