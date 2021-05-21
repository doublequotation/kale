package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/muesli/termenv"
)

func Command(cmd []string, verb string) int {
	command := exec.Command(cmd[0], cmd[1:]...)
	c := InitColors()
	stderr, err := command.StderrPipe()

	if err != nil {
		FPrint(c.Red, "Error", err)
		os.Exit(1)
	}

	command.Start()

	scanner := bufio.NewScanner(stderr)
	if scanner.Scan() == true {
		FPrint(c.Red, "Error", verb+" errors(s):")
		fmt.Println(termenv.String("-").Foreground(c.Red).Bold(), scanner.Text())
		for scanner.Scan() {
			fmt.Println(termenv.String("-").Foreground(c.Red).Bold(), scanner.Text())
		}
	}
	command.Wait()
	return command.Process.Pid
}
