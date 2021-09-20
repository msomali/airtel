package cli

import (
	"fmt"
	"github.com/techcraftlabs/airtel/api"
	"github.com/techcraftlabs/airtel/api/http"
	"github.com/techcraftlabs/airtel/internal/io"
	clix "github.com/urfave/cli/v2"
	"log"
	"os"
)

type (
	App struct {
		client api.Service
		app    *clix.App
	}
)

func add(svc api.Service) *clix.Command {
	return &clix.Command{
		Name:        "add",
		Aliases:     nil,
		Usage:       "add two numbers",
		UsageText:   "add [number] [number]",
		Description: "perform addition",
		ArgsUsage:   "args usage add command",
		Before:      beforeAddFunc,
		After:       afterAddFunc,
		Action: func(context *clix.Context) error {
			res, err := svc.Add(10, 20)
			if err != nil {
				return err
			}

			fmt.Printf("answer: %v\n", res)
			return nil
		},
		OnUsageError: onAddErrFunc,
	}
}

func div(svc api.Service) *clix.Command {
	return &clix.Command{
		Name:        "div",
		Aliases:     nil,
		Usage:       "divide two numbers",
		UsageText:   "div [number] [number]",
		Description: "perform division",
		ArgsUsage:   "args usage div command",
		Before:      beforeDivFunc,
		After:       afterDivFunc,
		Action: func(context *clix.Context) error {
			res, err := svc.Divide(100, 0)
			if err != nil {
				return err
			}
			fmt.Printf("answer: %v\n", res)
			return nil
		},
		OnUsageError: onDivErrFunc,
	}
}

func afterDivFunc(context *clix.Context) error {
	log.Printf("after div func")
	return nil
}

func beforeDivFunc(context *clix.Context) error {
	log.Printf("starting div func")
	return nil
}

func onDivErrFunc(context *clix.Context, err error, subcommand bool) error {
	log.Printf("error occurred during division: %v\n", err)
	return nil
}

func onAddErrFunc(context *clix.Context, err error, subcommand bool) error {
	fmt.Printf("add error: %v\n", err)
	return nil
}

func onAddActionFunc(context *clix.Context) error {
	return nil
}

func afterAddFunc(context *clix.Context) error {
	return nil
}

func beforeAddFunc(context *clix.Context) error {
	return nil
}

func commands(comm ...*clix.Command) []*clix.Command {
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

func New(base string, port uint64) *App {
	client := http.NewClient(base, port)
	author1 := &clix.Author{
		Name:  "Pius Alfred",
		Email: "me.pius1102@gmail.com",
	}

	app := &clix.App{
		Name:                 "calc",
		Usage:                "perform simple calculations",
		UsageText:            "calc add <first-integer> <second-integer>",
		Version:              "1.0.0",
		Description:          "perform simple addition and division of base64 integers",
		Commands:             commands(add(client), div(client)),
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
		client: client,
		app:    app,
	}
}

func beforeActionFunc(context *clix.Context) error {
	return nil
}

func afterActionFunc(context *clix.Context) error {
	return nil
}

func onActionFunc(context *clix.Context) error {
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
