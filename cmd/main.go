package main

import (
	"fmt"
	"kale/config"
	"kale/kl"
	"kale/utils"
	"os"
	"path/filepath"

	"github.com/akamensky/argparse"
	"github.com/common-nighthawk/go-figure"
	"github.com/muesli/termenv"
)

func logo(color termenv.Color) {
	rawArt := figure.NewFigure("kale", "isometric2", true)
	fig := fmt.Sprint(rawArt)
	logo := termenv.String(fig).Foreground(color).Bold()
	fmt.Println(logo)
}
func err(e error) {
	if e != nil {
		c := utils.InitColors()
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), e)
		os.Exit(0)
	}
}
func main() {
	c := utils.InitColors()
	logo(c.Green)

	root, _ := filepath.Abs(".")
	outFile := kl.Collect(root)
	Config := &config.Main{}
	kl.Run(outFile, Config)

	// Create new parser object
	parser := argparse.NewParser("kale", "A build system")
	buildCmd := parser.NewCommand("build", "Will start building workspace")
	versionCmd := parser.NewCommand("version", "shows current version")
	helpCmd := parser.NewCommand("doc", "Shows help information about specific topics")
	docSupportCmd := helpCmd.NewCommand("support", "Shows list of supported languages")
	// Parse input
	argErr := parser.Parse(os.Args)
	if argErr != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(""))
		os.Exit(0)
	}

	if docSupportCmd.Happened() {
		fmt.Println(termenv.String("Info:").Foreground(c.Cyan).Bold(), "Languages supported")
		fmt.Println(termenv.String("\t-").Foreground(c.Cyan).Bold(), "golang")
		fmt.Println(termenv.String("\t-").Foreground(c.Cyan).Bold(), "cpp/c")
	} else if versionCmd.Happened() {
		fmt.Println(termenv.String("1.1a").Foreground(c.Yellow))
	} else if buildCmd.Happened() {
		//fmt.Println(con)
		config.Do(Config)
	} else {
		fmt.Println(parser.Help(""))
	}
}
