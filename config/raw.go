package config

import (
	"fmt"
	Cmd "kale/commands"
	"kale/utils"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/muesli/termenv"
)

type fn func([]string)
type BuildSteps struct {
	Body [][]string
}

type Project struct {
	Output    string
	Sources   []string
	Extension string
	Params    []string
	Target    [][]string
}
type Config struct {
	Proj  Project    `toml:"project"`
	Steps BuildSteps `toml:"steps"`
}

func Read(content string) Config {
	var Conf Config
	if _, err := toml.Decode(content, &Conf); err != nil {
		c := utils.InitColors()
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), err)
		os.Exit(0)
	}
	return Conf
}
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

var buildConfig Config

func buildStep() {
	c := utils.InitColors()
	start := time.Now()
	conf := buildConfig
	if len(conf.Proj.Sources) == 0 || conf.Proj.Extension == "" || conf.Proj.Output == "" {
		c := utils.InitColors()
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), "Missing sources data, extension data, or output data in .KALE file. ")
		os.Exit(0)
	} else {
		var cmd []string = []string{}
		ext := conf.Proj.Extension
		if ext != "golang" {
			fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), "Uknown extension "+ext)
			fmt.Println(termenv.String("Info: ").Foreground(c.Cyan).Bold(), "Make sure it is a supported extension:")
			fmt.Println(termenv.String("\t- ").Foreground(c.Cyan).Bold(), "golang")
			fmt.Println(termenv.String("\t- ").Foreground(c.Cyan).Bold(), "cpp (coming soon)")
			os.Exit(0)
		} else {
			params := conf.Proj.Params
			cmd = append(cmd, "build", "-o", conf.Proj.Output)
			cmd = append(cmd, params...)
			Cmd.Build(cmd)
		}
	}
	duration := time.Since(start)
	fmt.Println(termenv.String("Time: ").Foreground(c.Cyan), duration.Seconds())
}
func doBuild(_ []string) {
	buildStep()
}
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Do(conf Config) {
	s := map[string][]string{
		"android":   {"arm"},
		"darwin":    {"386", "amd64", "arm", "arm64"},
		"dragonfly": {"amd64"},
		"freebsd":   {"386", "amd64", "arm"},
		"linux":     {"386", "amd64", "arm", "arm64", "ppc64", "ppc64le", "mips", "mipsle", "mips64", "mips64le"},
		"netbsd":    {"386", "amd64", "arm"},
		"openbsd":   {"386", "amd64", "arm"},
		"plan9":     {"386", "amd64"},
		"solaris":   {"amd64"},
		"windows":   {"386", "amd64"},
	}
	//osList := []string{"android", "linux", "darwin", "dragonfly", "freebsd", "linux", "netbsd", "openbsd", "plan9", "solaris", "windows"}
	//	archList := []string{"arm", "386", "amd64", "ppc64ppc64", "ppc64le", "mips", "mipsle", "mips64", "mips64le"}
	buildConfig = conf
	m := map[string]fn{
		"mkdir":  makeDir,
		"rmdir":  rmDir,
		"rmfile": rmFile,
		"mvfile": mvFile,
		"build":  doBuild,
	}
	c := utils.InitColors()
	for _, pair := range conf.Proj.Target {
		if len(s[pair[0]]) == 0 {
			fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), "Could not find build target operating system: "+pair[0])
			os.Exit(0)
		}
		if contains(s[pair[0]], pair[1]) == false {
			fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), pair[0]+" does not have architecure: "+pair[1])
			os.Exit(0)
		}
	}

	if len(conf.Steps.Body) != 0 {
		for _, set := range conf.Steps.Body {
			//if len(set) > 1 {
			operands := set[1:]
			if m[set[0]] == nil {
				fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), "Request "+set[0]+" does not exist.")
				os.Exit(0)
			}
			m[set[0]](operands)
			//}
		}
	}
}
