package command

import (
	"kale/utils"
)

type Builder struct {
	ProcName string
	Args     []string
	Cmd      string
	Output   string
	Target   []string
	Pid      int
	Env      []string
}

func (c *Builder) AddArgs(body ...string) {
	for _, p := range body {
		c.Args = append(c.Args, p)
	}
}

func (c *Builder) AddTarget(body ...string) {
	for _, p := range body {
		c.Target = append(c.Target, p)
	}
}

func (c *Builder) Construct() {
	cmd := []string{}
	if len(c.Env) > 0 {
		cmd = []string{"env", c.Env[0], c.Env[1], c.Cmd}
	} else {
		cmd = []string{c.Cmd}
	}
	c.AddArgs("-o", c.Output)
	cmd = append(cmd, c.Args...)
	for _, target := range c.Target {
		cmd = append(cmd, target)
	}
	c.Pid = utils.Command(cmd, c.ProcName)
}
