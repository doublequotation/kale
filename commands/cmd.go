package command

import (
	"fmt"
	"kale/utils"
	"strings"

	"github.com/muesli/termenv"
)

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

type GO struct {
	Platforms []string
	Params    []string
	Target    []string
	Out       string
}

func (g *GO) Build() {
	c := utils.InitColors()
	if len(g.Platforms) != 0 {
		for _, str := range g.Platforms {
			outName := g.Out
			g.Out = outName

			platDouble := strings.Split(str, " ")
			os := strings.Split(platDouble[0], "=")[1]
			arch := strings.Split(platDouble[0], "=")[1]

			//filePather := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
			//name := (filePather.FindStringSubmatch(g.Out))[2]
			g.Out = g.Out + "/" + os + "/" + "kale-" + arch

			build := Builder{ProcName: "Building", Output: g.Out, Cmd: "go", Env: platDouble}
			build.AddArgs("build")
			build.AddArgs(g.Params...)
			build.AddTarget(g.Target...)
			build.Construct()

			utils.FPrint(c.Green, "Compiled", g.Out)
			fmt.Println(termenv.String("\t- OS:").Foreground(c.Cyan).Bold(), g.Out)
			fmt.Println(termenv.String("\t- ARCHITECTURE:").Foreground(c.Cyan).Bold(), g.Out)
			fmt.Println(termenv.String("\t- PARAMS:").Foreground(c.Cyan).Bold(), strings.Join(g.Params, " "))
			g.Out = outName
		}
	} else {
		build := Builder{ProcName: "Building", Output: g.Out, Cmd: "go"}
		build.AddArgs("build")
		build.AddTarget(g.Target...)
		build.Construct()
		utils.FPrint(c.Green, "Compiled", g.Out)
	}
}
