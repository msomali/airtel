package cli

import (
	"fmt"
	"github.com/techcraftlabs/airtel/api"
	"github.com/techcraftlabs/airtel/api/http"
	"github.com/techcraftlabs/airtel/cli/internal"
	"github.com/techcraftlabs/airtel/internal/io"
	clix "github.com/urfave/cli/v2"
	"os"
)

type (
	App struct {
		client api.Service
		app    *clix.App
	}
)

func commands(client *http.Client) []*clix.Command {
	tokenCmd := internal.TokenCommand(client).Command()
	pushCmd := internal.PushCommand(client).Command()
	return appendCommands(tokenCmd,pushCmd)
}

func appendCommands(comm ...*clix.Command) []*clix.Command {
	var commands []*clix.Command
	for _, command := range comm {
		commands = append(commands, command)
	}
	return commands
}

func flags(fs ...clix.Flag) []clix.Flag {
	var flgs []clix.Flag
	for _, flg := range fs {
		flgs = append(flgs, flg)
	}
	return flgs
}

func authors(auth ...*clix.Author) []*clix.Author {
	var authors []*clix.Author
	for _, author := range auth {
		authors = append(authors, author)
	}
	return authors
}

func New(httpApiClient *http.Client ) *App {

	author1 := &clix.Author{
		Name:  "Pius Alfred",
		Email: "me.pius1102@gmail.com",
	}

	app := &clix.App{
		Name:                 "airtel",
		Usage:                "commandline tool to test/interact with Airtel Money API",
		UsageText:            "airtel push|disburse|balance|user|summary|refund|token|enquiry",
		Version:              "1.0.0",
		Description:          "interact with airtel money api",
		Commands:             commands(httpApiClient),
		Flags:                flags(),
		EnableBashCompletion: true,
		Before:               beforeActionFunc,
		After:                afterActionFunc,
		CommandNotFound:      onCommand404,
		OnUsageError:         onErrFunc,
		Authors:              authors(author1),
		Copyright:            "MIT Licence, Creative Commons",
		ErrWriter:            os.Stderr,
	}

	return &App{
		client: httpApiClient,
		app:    app,
	}
}

func beforeActionFunc(context *clix.Context) error {
	return nil
}

func afterActionFunc(context *clix.Context) error {
	return nil
}

func onCommand404(context *clix.Context, s string) {
	_, _ = fmt.Fprintf(io.Stderr, "not found: %s\n", s)
}

func onErrFunc(context *clix.Context, err error, subcommand bool) error {
	return nil
}

func (a *App) Run(args []string) error {
	return a.app.Run(args)
}
