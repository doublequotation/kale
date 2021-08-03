package kl

import (
	"encoding/json"
	"kale/config"
	"kale/utils"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/muesli/termenv"
	"google.golang.org/protobuf/proto"
	"rogchap.com/v8go"
)

type Path struct {
	Namespace string
	Property  string
	Named     bool
}

func parsePath(sPath string) Path {
	m, _ := regexp.MatchString(`^(\w+\/\/)?\w+`, sPath)
	path := Path{}
	sPath = strings.TrimSpace(sPath)
	if m == false {
		sPath = strings.Replace(sPath, "//", termenv.String("//").Foreground(c.Yellow).String(), -1)
		utils.FPrint(c.Red, "Error", "Invalid path", path)
		os.Exit(0)
	}
	if strings.Index(sPath, "//") != -1 {
		//result := make(map[string]interface{})
		//mapstructure.Decode(con, &result)
		//result[path]
		parts := strings.Split(sPath, "//")
		path.Namespace = strcase.ToCamel(strings.TrimSpace(parts[0]))
		if len(parts) == 1 {
			sPath = strings.Replace(sPath, "//", termenv.String("//").Foreground(c.Yellow).String(), -1)
			utils.FPrint(c.Red, "Error", "Invalid path", path)
			os.Exit(0)
		}
		path.Property = strcase.ToCamel(strings.TrimSpace(parts[1]))
		path.Named = true
	} else {
		path.Namespace = strcase.ToCamel(strings.TrimSpace(sPath))
		path.Property = strcase.ToCamel(strings.TrimSpace(sPath))
	}
	return path
}

func constructTuple(arch, os string) string {
	if os == "@android" && arch == "@android" {
		return "aarch64-unknown-linux-android"
	}

	return arch + "-unknown-" + os
}

func Run(path string, Config *config.Main) {
	iso, _ := v8go.NewIsolate()
	targetfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		t := &config.Target{}
		if len(info.Args()) < 2 {
			utils.FPrint(c.Red, "Error", "Function target requires 2 arguments.")
			os.Exit(0)
		}
		namePath := (info.Args()[0]).String()
		t.Name = &namePath
		if (info.Args()[1]).IsObject() == false {
			utils.FPrint(c.Red, "Error", "Function target requires 2nd argument to be an object.")
			os.Exit(0)
		}
		targetBytes, _ := (info.Args()[1]).Object().MarshalJSON()
		tempTarget := &config.Target{}
		err := json.Unmarshal(targetBytes, tempTarget)
		if err != nil {
			utils.FPrint(c.Red, "Error", err)
		}
		doubles := []string{}
		lastDouble := Path{}
		for _, a := range tempTarget.Tuples {
			tuple := parsePath(a)
			if tuple.Named == false && lastDouble.Property == "" {
				utils.FPrint(c.Red, "Error", "Expected path string in target definition named: "+namePath)
				os.Exit(0)
			} else if tuple.Named == false {
				tuple.Namespace = lastDouble.Namespace
			} else {
				lastDouble = tuple
			}

			doubles = append(doubles, strings.ToLower(constructTuple(tuple.Namespace, tuple.Property)))
		}
		t.Tuples = doubles
		Config.Platforms = append(Config.Platforms, t)
		return nil
	})
	outfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		t := &config.Target{}
		if len(info.Args()) < 2 {
			utils.FPrint(c.Red, "Error", "Function output requires 2 arguments.")
			os.Exit(0)
		}
		namePath := (info.Args()[0]).String()
		t.Name = &namePath
		if (info.Args()[1]).IsObject() == false {
			utils.FPrint(c.Red, "Error", "Function output requires 2nd argument to be an object.")
			os.Exit(0)
		}
		tempOut := &config.Output{}
		outBytes, _ := (info.Args()[1]).MarshalJSON()
		json.Unmarshal(outBytes, tempOut)
		targets := []string{}
		for _, t := range tempOut.Targets {
			path := parsePath(t)
			// appending indexes of targets into array as strings
			// this will be used to find the target when building sources
			index := sort.Search(len(Config.Platforms), func(i int) bool {
				return strings.ToLower(*Config.Platforms[i].Name) == strings.ToLower(path.Property)
			})
			targets = append(targets, strconv.Itoa(index))
		}
		//tempOut.Extension.String() != "golang" || tempOut.Extension.String() != "cpp" || tempOut.Extension.String() != "c" {
		// redefining the output targets with the indexes array
		tempOut.Targets = targets

		Config.Outs = append(Config.Outs, tempOut)

		return nil
	})
	runfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		cmd := []string{}
		for _, arg := range info.Args() {
			cmd = append(cmd, arg.String())
		}
		pid := utils.Command(cmd, "Scripting")
		v, _ := v8go.NewValue(iso, pid)
		return v
	})

	global, _ := v8go.NewObjectTemplate(iso)
	global.Set("target", targetfn)
	global.Set("output", outfn)
	global.Set("run_cmd", runfn)
	ctx, _ := v8go.NewContext(iso, global)
	source, _ := os.ReadFile(path)
	_, err := ctx.RunScript(string(source), "out.kl")
	if err != nil {
		utils.FPrint(c.Red, "Error", err)
	}
	encoded, encErr := proto.Marshal(Config)
	if encErr != nil {
		utils.FPrint(c.Red, "Error", "Failed to encode config:", encErr)
	}
	Ferr := os.WriteFile("conf", encoded, 0644)
	if encErr != nil {
		utils.FPrint(c.Red, "Error", "Failed to store config:", Ferr)
	}
}
