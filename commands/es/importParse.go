package main

import (
	"fmt"
	"os"
	"regexp"

	"rogchap.com/v8go"
)

func get(script string) []string {
	reg, _ := regexp.Compile(`import\s+(\")?(.*)+(\")?\s*(from\s+\"(.?)+\")?`)
	return reg.FindAllString(script, -1)
}

type Import struct {
	Filename string
	Contents string
}

func main() {
	if len(os.Args) >= 1 {
		filename := os.Args[0]
		content, _ := os.ReadFile(filename)
		iso, _ := v8go.NewIsolate()
		imports := []Import{}
		printfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			output := ""
			for _, arg := range info.Args() {
				output += arg.String() + " "
			}
			fmt.Printf("%v", output)
			return nil
		})
		vmfn, _ := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			name := info.Args()[1]
			imports = append(imports, Import{Filename: name.String(), Contents: script.String()})
			//ctx.RunScript(script.String(), name.String())
			return nil
		})
		global, _ := v8go.NewObjectTemplate(iso)
		global.Set("vm", vmfn)
		global.Set("print", printfn)
		ctx, _ := v8go.NewContext(iso, global)
		newScript := ""
		for _, im := range imports {
			newScript += im.Contents
			//ctx.RunScript(im.Contents, im.Filename)
		}
		//_, err := ctx.RunScript("print('foo', 'wow')", "print.js")
		//if err != nil {
		//	log.Fatalln(err)
		//}
	}
}
