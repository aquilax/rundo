package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/urfave/cli/v2"
)

const defaultFileName = ".rundo"
const nameSeparator = "="
const commentPrefix = "#"

type Command struct {
	Name        string
	Command     string
	Description string
}

func readCommands(r io.ReadCloser) ([]Command, error) {
	result := make([]Command, 0)
	s := bufio.NewScanner(r)
	var description string
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, commentPrefix) {
			// skip comment
			description = strings.TrimSpace(line[1:])
			continue
		}
		if strings.TrimSpace(line) == "" {
			// skip empty lines
			description = ""
			continue
		}
		index := strings.Index(line, nameSeparator)
		name := strings.TrimSpace(line[:index])
		command := strings.TrimSpace(line[index+1:])
		if index > -1 {
			result = append(result, Command{name, command, description})
		}
		description = ""
	}
	return result, nil
}

func printCommands(c []Command) error {
	fmt.Println("Available commands:")
	for _, c := range c {
		if c.Description != "" {
			fmt.Printf("# %s\n", c.Description)
		}
		fmt.Printf("%s = %s\n", c.Name, c.Command)
	}
	return nil
}

func runCommand(name, command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func getCommand(command string, cm []Command) (*Command, error) {
	for i := range cm {
		if cm[i].Name == command {
			return &cm[i], nil
		}
	}
	return nil, fmt.Errorf("Command %s not found", command)
}

func loadFile(fileName string) (io.ReadCloser, error) {
	f, err := os.Open(fileName)
	if err != nil {
		defer f.Close()
	}
	return f, err
}

func main() {
	// rundo command arg1 arg2

	app := &cli.App{
		Name:  "rundo",
		Usage: "Local scripts runner",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "Shows list of the defined commands in the current context",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "Dry run the command",
					},
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "Filename to use for commands",
					},
				},
				Action: func(c *cli.Context) error {
					fileName := "./" + defaultFileName
					if c.IsSet("file") {
						fileName = c.String("file")
					}
					f, err := loadFile(fileName)
					if err != nil {
						return err
					}

					cm, err := readCommands(f)
					if err != nil {
						return err
					}

					return printCommands(cm)
				},
			},
			{
				Name:      "run",
				Usage:     "Runs one of the commands",
				ArgsUsage: "[NAME] [ARGS]...",
				Aliases:   []string{"r"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "Dry run the command",
					},
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "Filename to use for commands",
					},
				},
				Action: func(c *cli.Context) error {
					// read commands
					fileName := "./" + defaultFileName
					if c.IsSet("file") {
						fileName = c.String("file")
					}
					f, err := loadFile(fileName)
					if err != nil {
						return err
					}

					cm, err := readCommands(f)
					if err != nil {
						return err
					}

					if c.NArg() == 0 {
						return printCommands(cm)
					}
					comm, err := getCommand(c.Args().First(), cm)
					if err != nil {
						return err
					}
					command, commandArgs, err := getCommandArgs(comm.Command)
					if err != nil {
						return err
					}
					args := append(commandArgs, c.Args().Tail()...)
					if c.IsSet("dry-run") {
						argsString := strings.Join(args, " ")
						fmt.Printf("# running alias `%s`: command: `%s` with arguments: `%s`\n", comm.Name, command, argsString)
						fmt.Printf("$ %s %s\n", command, argsString)
						return nil
					}
					return runCommand(comm.Name, command, args)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		os.Exit(1)
	}
}

func getCommandArgs(command string) (string, []string, error) {
	commandList, err := shellwords.Parse(command)
	return commandList[0], commandList[1:], err
}
