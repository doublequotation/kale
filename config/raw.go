package config

import (
	"fmt"
	command "kale/commands"
	"kale/utils"
	"os"
	"regexp"
	"strings"
	"time"

	cBuild "kale/commands/c"
	cppBuild "kale/commands/cpp"
	goBuild "kale/commands/golang"
	zapBuild "kale/commands/zap"

	"github.com/muesli/termenv"
)

type fn func([]string)

type Project struct {
	Output    string
	Extension string
	Sources   []string
	Params    []string
	Target    [][]string
}
type Zap struct {
	Env     []string
	Sources []string
}

type Config struct {
	Proj Project
	Zap  Zap
}

func TupleToDouble(Tuple string) []string {
	return strings.Split(Tuple, "-unknown-")
}

var buildConfig Config
var validPairs [][]string = [][]string{}

func pairToEnv() []string {
	pgroup := []string{}
	for _, p := range validPairs {
		pgroup = append(pgroup, "GOOS="+p[0]+" GOARCH="+p[1])
	}
	return pgroup
}
func ZapStep() {
	conf := buildConfig
	if len(conf.Zap.Env) == 0 && len(conf.Zap.Sources) >= 1 {
		c := utils.InitColors()
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), "Cannot inject environment variables with no variables.")
		os.Exit(0)
	}
	listing := []string{}
	for _, env := range conf.Zap.Env {
		Var, exists := os.LookupEnv(env)
		if exists == false {
			c := utils.InitColors()
			fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), "No environment variable named", termenv.String(env).Foreground(c.Cyan))
			os.Exit(0)
		}
		listing = append(listing, env+"="+Var)
	}

	for _, name := range conf.Zap.Sources {
		zapBuild.Zap(listing, name)
	}
}
func buildStep() {
	c := utils.InitColors()
	start := time.Now()
	conf := buildConfig
	if len(conf.Proj.Sources) == 0 || conf.Proj.Extension == "" || conf.Proj.Output == "" {
		c := utils.InitColors()
		fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), "Missing sources data, extension data, or output data in .KALE file. ")
		os.Exit(0)
	} else {
		ext := conf.Proj.Extension
		if ext == "golang" {
			ZapStep()
			params := conf.Proj.Params
			cmd := goBuild.GO{}
			if strings.HasSuffix(conf.Proj.Sources[0], "/*") {
				path := strings.Replace(conf.Proj.Sources[0], "/*", "", 1)
				files, dirErr := os.ReadDir(path)
				if dirErr != nil {
					c := utils.InitColors()
					fmt.Println(termenv.String("Error: ").Foreground(c.Red).Bold(), dirErr)
					os.Exit(0)
				}
				conf.Proj.Sources = []string{}
				for _, dir := range files {
					if strings.HasSuffix(dir.Name(), ".go") {
						filePather := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
						cmd.Out = (filePather.FindStringSubmatch(dir.Name()))[2]
						cmd.Params = append(cmd.Params, params...)
						cmd.Target = append(cmd.Target, dir.Name())
						cmd.Platforms = pairToEnv()
						cmd.Build()
						conf.Proj.Sources = append(conf.Proj.Sources, conf.Proj.Output)
					}
				}
			} else if len(conf.Proj.Sources) > 1 {
				for _, file := range conf.Proj.Sources {
					filePather := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
					cmd.Out = (filePather.FindStringSubmatch(file))[2]
					cmd.Params = append(cmd.Params, params...)
					cmd.Target = append(cmd.Target, file)
					cmd.Platforms = pairToEnv()
					cmd.Build()
				}
			} else {
				cmd.Target = append(cmd.Target, conf.Proj.Sources[0])
				cmd.Platforms = pairToEnv()
				cmd.Out = conf.Proj.Output
				cmd.Build()
			}
		} else if ext == "cpp" || ext == "c" {
			home, _ := os.UserHomeDir()
			os.Setenv("WORKDIR", home+"/.config/kale")
			os.MkdirAll(home+"/.config/kale", 0755)
			if strings.HasSuffix(conf.Proj.Sources[0], "/*") {
				path := strings.Replace(conf.Proj.Sources[0], "/*", "", 1)
				files, dirErr := os.ReadDir(path)
				if dirErr != nil {
					c := utils.InitColors()
					fmt.Println(termenv.String("Error:").Foreground(c.Red).Bold(), dirErr)
					os.Exit(0)
				}
				conf.Proj.Sources = []string{}
				var cmd []command.Builder = []command.Builder{}
				//if buildConfig.C.Compiler == "" {
				//fmt.Println(termenv.String("Info:").Foreground(c.Cyan).Bold(), "Defaulting to g++/gcc")
				//}
				objects := []string{}
				for _, dir := range files {
					if strings.HasSuffix(dir.Name(), "cpp") || strings.HasSuffix(dir.Name(), "cc") {
						filePather := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
						name := (filePather.FindStringSubmatch(dir.Name()))[2]
						pre := command.Builder{Cmd: "g++", Output: os.Getenv("WORKDIR") + "/" + name + ".i", Target: []string{path + "/" + dir.Name()}}
						pre.AddArgs("-E")
						pre.AddArgs(conf.Proj.Params...)
						pre.AddTarget(path + "/" + dir.Name())

						asm := command.Builder{Cmd: "g++", Output: os.Getenv("WORKDIR") + "/" + name + ".S", Target: []string{path + "/" + dir.Name()}}
						asm.AddArgs("-S")
						asm.AddArgs(conf.Proj.Params...)
						asm.AddTarget(path + "/" + dir.Name())

						obj := command.Builder{Cmd: "g++", Output: os.Getenv("WORKDIR") + "/" + name + ".o", Target: []string{path + "/" + dir.Name()}}
						obj.AddArgs("-c")
						obj.AddArgs(conf.Proj.Params...)
						obj.AddTarget(path + "/" + dir.Name())

						objects = append(objects, os.Getenv("WORKDIR")+"/"+name+".o")
					} else if strings.HasSuffix(dir.Name(), "c") {
						filePather := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
						name := (filePather.FindStringSubmatch(dir.Name()))[2]
						pre := command.Builder{Cmd: "gcc", Output: os.Getenv("WORKDIR") + "/" + name + ".i", Target: []string{path + "/" + dir.Name()}}
						pre.AddArgs("-E")
						pre.AddArgs(conf.Proj.Params...)
						pre.AddTarget(path + "/" + dir.Name())

						asm := command.Builder{Cmd: "gcc", Output: os.Getenv("WORKDIR") + "/" + name + ".s", Target: []string{path + "/" + dir.Name()}}
						asm.AddArgs("-S")
						asm.AddArgs(conf.Proj.Params...)
						asm.AddTarget(path + "/" + dir.Name())

						obj := command.Builder{Cmd: "gcc", Output: os.Getenv("WORKDIR") + "/" + name + ".o", Target: []string{path + "/" + dir.Name()}}
						obj.AddArgs("-c")
						obj.AddArgs(conf.Proj.Params...)
						obj.AddTarget(path + "/" + dir.Name())

						cmd = append(cmd, pre, asm, obj)

						objects = append(objects, os.Getenv("WORKDIR")+"/"+name+".o")
					}
				}
				if ext == "cpp" {
					builder := cppBuild.CPP{Args: conf.Proj.Params, Steps: cmd, Objects: objects, Out: conf.Proj.Output}
					builder.CppBuild()
				} else {
					builder := cBuild.C{Args: conf.Proj.Params, Steps: cmd, Objects: objects, Out: conf.Proj.Output}
					builder.CBuild()
					//cBuild.CBuild(cmd, conf.Proj.Output, objects, conf.Proj.Params)
				}
				os.RemoveAll(home + "/.config/kale")
			}
		} else {
			fmt.Println(termenv.String("Error:").Foreground(c.Red).Bold(), "Uknown extension "+ext)
			fmt.Println(termenv.String("Info:").Foreground(c.Cyan).Bold(), "Make sure it is a supported extension:")
			fmt.Println(termenv.String("\t-").Foreground(c.Cyan).Bold(), "golang")
			fmt.Println(termenv.String("\t-").Foreground(c.Cyan).Bold(), "cpp/c")
			os.Exit(0)
		}
	}
	duration := time.Since(start)
	for _, path := range conf.Zap.Sources {
		zapBuild.Transfer(path)
	}
	fmt.Println(termenv.String("Time:").Foreground(c.Cyan), duration.Seconds())
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

func Do(Conf *Main) {
	c := utils.InitColors()
	//	methods := map[string]fn{
	//		"mkdir":  makeDir,
	//		"rmdir":  rmDir,
	//		"rmfile": rmFile,
	//		"mvfile": mvFile,
	//		"copy":   cpAny,
	//		"build":  doBuild,
	//	}
	for _, conf := range Conf.Outs {
		if conf.Extension.String() == "cpp" {
			if len(conf.Targets) != 0 {
				utils.FPrint(c.Red, "Error", "Cpp does not support multiple build targets currently.")
				utils.FPrint(c.Cyan, "Info", "This will be implemented later.")
				fmt.Println(termenv.String("\t-").Foreground(c.Cyan).Bold(), "If you want to implement a feature, contribute to this project: ", termenv.String("https://github.com/doublequotation/kale").Foreground(c.Yellow))
			}
		} else if conf.Extension.String() == "golang" {
			s := map[string][]string{
				"android":   {"arm"},
				"darwin":    {"386", "amd64", "arm64"},
				"dragonfly": {"amd64"},
				"freebsd":   {"386", "amd64", "arm"},
				"linux":     {"386", "amd64", "arm", "arm64", "ppc64", "ppc64le", "mips", "mipsle", "mips64", "mips64le"},
				"netbsd":    {"386", "amd64", "arm"},
				"openbsd":   {"386", "amd64", "arm"},
				"plan9":     {"386", "amd64"},
				"solaris":   {"amd64"},
				"windows":   {"386", "amd64"},
			}
			for _, TupleTarget := range conf.Targets {
				pair := TupleToDouble(TupleTarget)
				if len(s[pair[0]]) == 0 {
					utils.FPrint(c.Red, "Error", "Could not find build target operating system: "+pair[0])
					os.Exit(0)
				}
				validPairs = append(validPairs, []string{pair[0]})
				if len(pair[1:]) > 1 {
					for i, target := range pair[1:] {
						if contains(s[pair[0]], target) == false {
							utils.FPrint(c.Red, "Error", pair[0]+" does not have architecure: "+target)
							os.Exit(0)
						}
						if i == 0 {
							validPairs[len(validPairs)-1] = append(validPairs[len(validPairs)-1], target)
						} else {
							validPairs = append(validPairs, []string{pair[0], target})
						}
					}
				} else {
					if contains(s[pair[0]], pair[1]) == false {
						utils.FPrint(c.Red, "Error", pair[0]+" does not have architecure: "+pair[1])
						os.Exit(0)
					}
					validPairs[len(validPairs)-1] = append(validPairs[len(validPairs)-1], pair[1])
				}
			}
		}
	}
	//if len(conf.Steps.Body) != 0 {
	//	for _, set := range conf.Steps.Body {
	//		operands := set[1:]
	//		if methods[set[0]] == nil {
	//			utils.FPrint(c.Red, "Error", "Request "+set[0]+" does not exist.")
	//			os.Exit(0)
	//		}
	///		methods[set[0]](operands)
	//	}
	//}
}
