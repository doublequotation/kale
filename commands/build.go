package command

import "kale/utils"

type Builder struct {
	ProcName string
	Args     []string
	Cmd      string
	Output   string
	Target   []string
	Pid      int
}

func (c *Builder) AddArgs(body ...string) {
	for _, p := range body {
		c.Args = append(c.Args, p)
	}
}

func (c *Builder) AddTarget(body ...string) {
	for _, p := range body {
		c.Args = append(c.Args, p)
	}
}

func (c *Builder) Construct() {
	c.AddArgs("-o", c.Output)
	cmd := []string{c.Cmd}
	cmd = append(cmd, c.Args...)
	for _, target := range c.Target {
		c.AddArgs(target)
	}

	c.Pid = utils.Command(cmd, c.ProcName)
}
