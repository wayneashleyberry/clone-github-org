package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	git "gopkg.in/src-d/go-git.v4"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	org := flag.String("o", "", "organisation")
	token := flag.String("t", "", "token")
	flag.Parse()

	if *org == "" {
		log.Fatal("missing -o flag")
	}

	if *token == "" {
		log.Fatal("missing -t flag")
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	repos, _, err := client.Repositories.ListByOrg(ctx, *org, &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 1000,
		},
	})
	must(err)

	err = os.MkdirAll(*org, os.ModePerm)
	must(err)

	for _, repo := range repos {
		url := fmt.Sprintf("https://%s:x-oauth-basic@github.com/%s", *token, *repo.FullName)

		_, err = git.PlainClone(*repo.FullName, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		if err != nil {
			fmt.Printf("%s: %s\n", *repo.CloneURL, err)
		}
	}
}
