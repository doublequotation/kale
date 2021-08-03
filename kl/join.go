package kl

import (
	"kale/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/evanw/esbuild/pkg/api"
)

var c = utils.InitColors()

func Collect(root string) string {
	rootDir, ferr := os.ReadDir(root)

	if os.IsNotExist(ferr) {
		utils.FPrint(c.Red, "Error", "Cannot read files in project root:", root)
	}
	names := []string{}
	imports := ""
	for _, f := range rootDir {
		if strings.HasSuffix(f.Name(), ".kl") && f.IsDir() == false {
			names = append(names, f.Name())
			p, _ := filepath.Abs(f.Name())
			imports += `import "` + p + "\"\n"
		}
	}
	home, _ := os.UserHomeDir()
	path := home + "/.config/kale/"
	err := os.WriteFile(path+"base.kl", []byte(imports), 0644)
	if err != nil {
		utils.FPrint(c.Red, "Error", err)
	}
	r := api.Build(api.BuildOptions{
		EntryPoints:       []string{path + "base.kl"},
		Outfile:           path + "out.kl",
		Write:             true,
		Bundle:            true,
		MinifyIdentifiers: true,
		Loader: map[string]api.Loader{
			".kl": api.LoaderJS,
		},
	})
	os.Remove(path + "base.kl")
	if len(r.Errors) > 0 {
		for _, err := range r.Errors {
			utils.FPrint(c.Red, "Error", err.Text)
			utils.FPrint(c.Cyan, "Info", err.Location.File+":"+strconv.Itoa(err.Location.Line))
			quick.Highlight(os.Stdout, err.Location.LineText, "javascript", "noop", "xcode")
		}
	}
	return path + "/out.kl"
}
