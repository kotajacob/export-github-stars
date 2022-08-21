package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v46/github"
	"github.com/muesli/reflow/wordwrap"
	"golang.org/x/oauth2"
)

func main() {
	flag.Parse()
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: flag.Arg(0)},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	i := 1
	for {
		// List all starred repos for the authenticated user.
		opts := github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
				Page:    i,
			},
		}

		stars, resp, err := client.Activity.ListStarred(ctx, "", &opts)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				log.Fatalln("hit rate limit")
			}
			log.Fatalf("failed listing stars: %v\n", err)
		}

		for _, star := range stars {
			repo := star.GetRepository()

			fmt.Println(repo.GetHTMLURL())

			desc := wordwrap.String(repo.GetDescription(), 80)
			if desc != "" {
				fmt.Println("Desc:", desc)
			}

			lang := repo.GetLanguage()
			if lang == "" {
				lang = "Unknown"
			}

			fmt.Println("Lang:", lang)
			fmt.Println("Stars:", repo.GetStargazersCount())

			if len(repo.Topics) != 0 {
				topics := strings.Join(repo.Topics, ", ")
				fmt.Println("Topics:", topics)
			}

			fmt.Println()
		}
		if resp.NextPage == 0 {
			break
		}
		i += 1
	}
}
