package kl

import (
	"fmt"
	"kale/config"
	"kale/utils"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/muesli/termenv"
	"rogchap.com/v8go"
)

type Path struct {
	Namespace string
	Property  string
}

func parsePath(sPath string) Path {
	m, _ := regexp.MatchString(`^(\w+\/\/)?\w+`, sPath)
	path := Path{}
	sPath = strcase.ToCamel(sPath)
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
		path.Namespace = strings.TrimSpace(parts[0])
		if len(parts) == 1 {
			sPath = strings.Replace(sPath, "//", termenv.String("//").Foreground(c.Yellow).String(), -1)
			utils.FPrint(c.Red, "Error", "Invalid path", path)
			os.Exit(0)
		}
		path.Property = strings.TrimSpace(parts[1])
	} else {
		path.Namespace = strings.TrimSpace(sPath)
		path.Property = strings.TrimSpace(sPath)
	}
	return path
}

func Run(path string) {
	con := config.Config{}
	iso, _ := v8go.NewIsolate()
	setfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		namePath := (info.Args()[0]).String()
		val := (info.Args()[1]).String()
		pPath := parsePath(namePath)
		r := reflect.ValueOf(con)
		f := reflect.Indirect(r).FieldByName(pPath.Namespace)
		f.SetString(val)
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
	global.Set("set_option", setfn)
	global.Set("run_cmd", runfn)
	ctx, _ := v8go.NewContext(iso, global)
	source, _ := os.ReadFile(path)
	fmt.Println("wow")
	_, err := ctx.RunScript(string(source), "out.kl")
	if err != nil {
		utils.FPrint(c.Red, "Error", err)
	}
}
