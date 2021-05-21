package config

import (
	"fmt"
	"kale/utils"
	"os"

	"github.com/muesli/termenv"
)

func makeDir(dirs []string) {
	for _, dir := range dirs {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			c := utils.InitColors()
			fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), err)
			os.Exit(0)
		}
	}
}

func copyFile(sourceFile, destinationFile string) {
	input, err := os.ReadFile(sourceFile)
	if err != nil {
		c := utils.InitColors()
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), err)
		os.Exit(0)
	}

	err = os.WriteFile(destinationFile, input, 0644)
	if err != nil {
		c := utils.InitColors()
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), err)
		os.Exit(0)
	}
}
func cpAny(dirs []string) {
	for _, dir := range dirs[1:] {
		copyFile(dirs[0], dir)
	}
}
func rmDir(dirs []string) {
	for _, dir := range dirs {
		err := os.RemoveAll(dir)
		if err != nil {
			c := utils.InitColors()
			fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), err)
			os.Exit(0)
		}
	}
}
func rmFile(dirs []string) {
	for _, dir := range dirs {
		err := os.Remove(dir)
		if err != nil {
			c := utils.InitColors()
			fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), err)
			os.Exit(0)
		}
	}
}
func mvFile(dirs []string) {
	if len(dirs) > 2 {
		c := utils.InitColors()
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), "Invalid move file request. More than 2 origin paths.")
		os.Exit(0)
	}
	err := os.Rename(dirs[0], dirs[1])
	if err != nil {
		c := utils.InitColors()
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), err)
		os.Exit(0)
	}
}
