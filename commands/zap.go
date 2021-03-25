package command

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var path string

type valid struct {
	Name  string
	Value string
}

func Zap(evs []string, name string) {
	var aval []valid
	for _, v := range evs {
		tv := strings.Split(v, "=")
		aval = append(aval, valid{Name: tv[0], Value: tv[1]})
	}
	// fmt.Println(flag.Param)
	// f, _ := os.OpenFile(flag.Param, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	contentB, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err)
	}
	er := ioutil.WriteFile(name+".copy", []byte(contentB), 0644)
	if er != nil {
		fmt.Println(er)
	}
	path = name
	content := string(contentB)
	sN := strings.Split(content, "\n")
	for i, line := range sN {
		if m, _ := regexp.MatchString(`//\s*@zap\s+var:`, line); m == true {
			if m, _ := regexp.MatchString(`var\s+\S+\s+\w+\s*=`, sN[i+1]); m == true {
				params := strings.Fields(sN[i+1])
				name := params[1]
				tp := params[2]
				r := regexp.MustCompile(`//\s*@zap\s+var:`)
				s := r.Split(line, -1)

				for _, a := range aval {
					if a.Name == strings.Replace(strings.Join(s[1:], ""), " ", "", -1) {
						newVar := fmt.Sprintln("var", name, tp, "= \""+a.Value+"\"")
						sN[i+1] = newVar
					}
				}
			}
		}
	}

	file := strings.Join(sN, "\n")
	ioutil.WriteFile(name, []byte(file), 0644)
	//build()
}
