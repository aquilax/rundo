package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

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

func getCommand(command string, cm []Command) (string, string, error) {
	for i := range cm {
		if cm[i].Name == command {
			return cm[i].Name, cm[i].Command, nil
			``
		}
	}
	return "", "", fmt.Errorf("Command %s not found", command)
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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "Filename to use for commands",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Dry run the command",
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

			cm, _ := readCommands(f)
			if c.NArg() == 0 {
				return printCommands(cm)
			}
			name, command, err := getCommand(c.Args().First(), cm)
			if err != nil {
				return err
			}
			commandList := strings.Fields(command)
			command = commandList[0]
			args := append(commandList[1:], c.Args().Tail()...)
			if c.IsSet("dry-run") {
				argsString := strings.Join(args, " ")
				fmt.Printf("# running alias `%s`: command: `%s` with arguments: `%s`\n", name, command, argsString)
				fmt.Printf("$ %s %s\n", command, argsString)
				return nil
			}
			return runCommand(name, command, args)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
