package cBuild

import (
	command "kale/commands"
)

type C struct {
	Steps   []command.Builder
	Out     string
	Objects []string
	Args    []string
}

func (c *C) CBuild() {
	for _, cmd := range c.Steps {
		cmd.Construct()
	}

	cmd := command.Builder{Cmd: "gcc"}
	cmd.AddSteps(c.Args...)
	cmd.AddTarget(c.Objects...)
	cmd.Construct()
}
