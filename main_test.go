package main

import (
	"io"
	"reflect"
	"testing"
)

func Test_readCommands(t *testing.T) {
	type args struct {
		r io.ReadCloser
	}
	tests := []struct {
		name    string
		args    args
		want    []Command
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readCommands(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("readCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_printCommands(t *testing.T) {
	type args struct {
		c []Command
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := printCommands(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("printCommands() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_runCommand(t *testing.T) {
	type args struct {
		name    string
		command string
		args    []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runCommand(tt.args.name, tt.args.command, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getCommand(t *testing.T) {
	type args struct {
		command string
		cm      []Command
	}
	tests := []struct {
		name    string
		args    args
		want    *Command
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCommand(tt.args.command, tt.args.cm)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadFile(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadFile(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_execute(t *testing.T) {
	tests := []struct {
		name        string
		fullCommand string
		want        error
	}{
		{
			"pipe",
			"bash -c 'ls -1 | sort'",
			nil,
		},
		{
			"simple",
			"ls",
			nil,
		},
		{
			"simple-args",
			"ls -lah",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command, commandArgs, _ := getCommandArgs(tt.fullCommand)
			got := runCommand(tt.name, command, commandArgs)
			if got != tt.want {
				t.Errorf("runCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
