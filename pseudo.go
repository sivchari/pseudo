package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/google/go-github/v69/github"
	"github.com/sivchari/commander"
	"github.com/spf13/pflag"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	cmd := commander.NewCommandManager().Build()
	client := github.NewClient(nil)
	cmd.Register(&pseudoCommand{client: client})
	if err := cmd.Run(ctx); err != nil {
		panic(err)
	}
}

type pseudoCommand struct {
	client *github.Client
	owner  string
	repo   string
	sha    string
}

var _ commander.Commander = (*pseudoCommand)(nil)

var ErrNoTags = errors.New("no tags found")

func (p *pseudoCommand) Run(ctx context.Context) error {
	if p.sha == "" {
		tags, _, err := p.client.Repositories.ListTags(ctx, p.owner, p.repo, nil)
		if err != nil {
			return err
		}
		if len(tags) == 0 {
			return ErrNoTags
		}
		fmt.Println(tags[0].GetName())
		return nil
	}
	commit, _, err := p.client.Git.GetCommit(ctx, p.owner, p.repo, p.sha)
	if err != nil {
		return err
	}
	sha := commit.GetSHA()
	date := commit.GetCommitter().GetDate()
	fmt.Printf("v0.0.0-%s-%s\n", date.Format("20060102150405"), sha[:12])
	return nil
}

func (p *pseudoCommand) Name() string {
	return "pseudo"
}

func (p *pseudoCommand) Short() string {
	return "pseudos generates pseudo version of the given URL"
}

func (p *pseudoCommand) Long() string {
	return "pseudos generates pseudo version of the given URL"
}

func (p *pseudoCommand) SetFlags(f *pflag.FlagSet) {
	f.StringVar(&p.owner, "owner", "", "owner of the repository")
	f.StringVar(&p.repo, "repo", "", "repository name")
	f.StringVar(&p.sha, "sha", "", "commit sha")
}
