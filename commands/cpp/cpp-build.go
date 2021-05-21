package cppBuild

import "kale/utils"

func CppBuild(steps [][]string, out string, objects []string, args []string) {
	for _, cmd := range steps {
		utils.Command(cmd, "Building")
	}

	cmd := []string{"g++"}
	cmd = append(cmd, args...)
	cmd = append(cmd, "-o", out)
	cmd = append(cmd, objects...)
	utils.Command(cmd, "Compiling")
}
