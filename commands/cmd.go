package command

import (
	"bufio"
	"fmt"
	"kale/utils"
	"os"
	"os/exec"
	"strings"

	"github.com/muesli/termenv"
)

func rawBuild(cmdStr ...string) {
	command := exec.Command(cmdStr[0], cmdStr[1:]...)
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
func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
func Build(params []string, pre []string) {
	c := utils.InitColors()
	if len(pre) != 0 {
		for _, str := range pre {
			outName := params[indexOf("-o", params)+1]
			params[indexOf("-o", params)+1] = outName
			// params := strings.Join(params, " ")
			nStr := strings.Split(str, " ")
			params[indexOf("-o", params)+1] = params[indexOf("-o", params)+1] + "-" + strings.Split(nStr[1], "=")[1] + "-" + strings.Split(nStr[2], "=")[1]
			nStr = append(nStr, params...)
			fmt.Println(termenv.String("Compiled: ").Foreground(c.Green).Bold(), params[indexOf("-o", params)+1])
			fmt.Println(termenv.String("\t- OS:  ").Foreground(c.Cyan).Bold(), params[indexOf("-o", params)+1])
			fmt.Println(termenv.String("\t- ARCHITECTURE:  ").Foreground(c.Cyan).Bold(), params[indexOf("-o", params)+1])
			fmt.Println(termenv.String("\t- PARAMS:  ").Foreground(c.Cyan).Bold(), strings.Join(params, " "))
			rawBuild(nStr...)
			params[indexOf("-o", params)+1] = outName
		}
	} else {
		rawBuild(params...)
	}
}
