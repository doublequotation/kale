package command

import (
	"bufio"
	"fmt"
	"kale/utils"
	"os"
	"os/exec"

	"github.com/muesli/termenv"
)

func Build(params []string) {
	command := exec.Command("go", params...)
	c := utils.InitColors()
	stderr, err := command.StderrPipe()

	if err != nil {
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), err)
		os.Exit(1)
	}

	command.Start()

	scanner := bufio.NewScanner(stderr)
	if scanner.Scan() == true {
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), "Compiling error(s):")
		fmt.Println(termenv.String("- ").Foreground(c.Red).Bold(), scanner.Text())
		for scanner.Scan() {
			fmt.Println(termenv.String("- ").Foreground(c.Red).Bold(), scanner.Text())
		}
	}
	command.Wait()
}
