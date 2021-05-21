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
func Build(params []string, pre []string, output string) {
	c := utils.InitColors()
	if len(pre) != 0 {
		for _, str := range pre {
			outName := params[indexOf("-o", params)+1]
			params[indexOf("-o", params)+1] = outName
			// params := strings.Join(params, " ")
			nStr := strings.Split(str, " ")
			//if _, err := os.Mkdir(strings.Split(nStr[1], "=")[1], 0755); os.IsNotExist(err) == true {
			//}
			params[indexOf("-o", params)+1] = output + strings.Split(nStr[1], "=")[1] + "/" + params[indexOf("-o", params)+1] + "-" + strings.Split(nStr[2], "=")[1]
			nStr = append(nStr, params...)
			utils.Command(nStr, "Building")
			utils.FPrint(c.Green, "Compiled", params[indexOf("-o", params)+1])
			fmt.Println(termenv.String("\t- OS:").Foreground(c.Cyan).Bold(), params[indexOf("-o", params)+1])
			fmt.Println(termenv.String("\t- ARCHITECTURE:").Foreground(c.Cyan).Bold(), params[indexOf("-o", params)+1])
			fmt.Println(termenv.String("\t- PARAMS:").Foreground(c.Cyan).Bold(), strings.Join(params[1:], " "))
			params[indexOf("-o", params)+1] = outName
		}
	} else {
		params[indexOf("-o", params)+1] = output + params[indexOf("-o", params)+1]
		utils.Command(params, "Building")
		utils.FPrint(c.Green, "Compiled", params[indexOf("-o", params)+1])
	}
}
