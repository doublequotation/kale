package main

import (
	"fmt"
	"regexp"

	"rogchap.com/v8go"
)

func get(script string) []string {
	reg, _ := regexp.Compile(`import\s+(\")?(.*)+(\")?\s*(from\s+\"(.?)+\")?`)
	return reg.FindAllString(script, -1)
}
func main() {
	iso, _ := v8go.NewIsolate()

	printfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		output := ""
		for _, arg := range info.Args() {
			output += arg.String() + " "
		}
		fmt.Printf("%v", output)
		return nil
	})
	global, _ := v8go.NewObjectTemplate(iso)
	global.Set("print", printfn)
	ctx, _ := v8go.NewContext(iso, global)
	ctx.RunScript("print('foo', 'wow')", "print.js")
}
