package cppBuild

import (
	command "kale/commands"
)

type CPP struct {
	Steps   []command.Builder
	Out     string
	Objects []string
	Args    []string
}

func (cpp *CPP) CppBuild() {
	for _, cmd := range cpp.Steps {
		cmd.Construct()
	}

	cmd := command.Builder{Cmd: "g++", Output: cpp.Out}
	cmd.AddArgs(cpp.Args...)
	cmd.AddTarget(cpp.Objects...)
	cmd.Construct()
}
