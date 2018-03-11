package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	gitHubToken  string
	creationOpts CreateOpts
)

type CreateOpts struct {
	Public      bool
	Name        string
	Description string
	Contents    string
}

type TokenSource oauth2.Token

func (t *TokenSource) Token() (*oauth2.Token, error) {
	return (*oauth2.Token)(t), nil
}

func (c CreateOpts) Valid() error {
	if c.Name == "" {
		return errors.New("A gist name is required")
	}

	return nil
}

type GistFiles map[github.GistFilename]github.GistFile

func createGist(options CreateOpts) error {
	ts := TokenSource{AccessToken: gitHubToken}
	client := github.NewClient(
		oauth2.NewClient(oauth2.NoContext, &ts),
	)

	files := make(GistFiles)
	name := github.GistFilename(options.Name)
	files[name] = github.GistFile{Content: &options.Contents}

	gist, _, err := client.Gists.Create(context.TODO(), &github.Gist{
		Files:       files,
		Public:      &options.Public,
		Description: &options.Description,
	})

	if err != nil {
		return err
	}

	fmt.Println(*gist.HTMLURL)

	return nil
}

func main() {
	// FIXME: Should probably use cobra or something for this when adding more commands.
	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	createCommand.BoolVar(&creationOpts.Public, "public", false, "Should the created gist be public")
	createCommand.StringVar(&creationOpts.Name, "name", "", "The name of the file")
	createCommand.StringVar(&creationOpts.Description, "desc", "", "The description of the gist, can be blank")
	createCommand.StringVar(&creationOpts.Contents, "content", "", "The content of the gist, can also be passed via stdin")
	createCommand.StringVar(&gitHubToken, "token", os.Getenv("GITHUB_TOKEN"), "The GitHub API token to use")

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage of %s <command> [options]:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  create\tcreate a new gist")
		fmt.Fprintf(os.Stderr, "\n\n")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "create":
		createCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(2)
	}

	if createCommand.Parsed() {
		if err := creationOpts.Valid(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			createCommand.Usage()
			fmt.Fprintf(os.Stderr, "\n\n")
			os.Exit(2)
		}

		if creationOpts.Contents == "" {
			data, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n\n", err)
				os.Exit(1)
			}
			creationOpts.Contents = string(data)
		}

		if err := createGist(creationOpts); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n\n", err)
			os.Exit(1)
		}
	}
}
