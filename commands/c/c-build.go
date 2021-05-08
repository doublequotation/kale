package cBuild

import (
	"bufio"
	"fmt"
	"kale/utils"
	"os"
	"os/exec"

	"github.com/muesli/termenv"
)

func Command(cmd []string, verb string) {
	command := exec.Command(cmd[0], cmd[1:]...)
	c := utils.InitColors()
	stderr, err := command.StderrPipe()

	if err != nil {
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), err)
		os.Exit(1)
	}

	command.Start()

	scanner := bufio.NewScanner(stderr)
	if scanner.Scan() == true {
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), verb+" error(s):")
		fmt.Println(termenv.String("- ").Foreground(c.Red).Bold(), scanner.Text())
		for scanner.Scan() {
			fmt.Println(termenv.String("- ").Foreground(c.Red).Bold(), scanner.Text())
		}
	}
	command.Wait()
}
func Build(steps [][]string, out string, objects []string) {
	for _, cmd := range steps {
		Command(cmd, "Building")
	}

	cmd := []string{"g++"}
	cmd = append(cmd, "-o", out)
	cmd = append(cmd, objects...)
	Command(cmd, "Compiling")
}
