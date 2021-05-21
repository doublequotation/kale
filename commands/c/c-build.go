package cBuild

import "kale/utils"

func CBuild(steps [][]string, out string, objects []string, args []string) {
	for _, cmd := range steps {
		utils.Command(cmd, "Building")
	}

	cmd := []string{"gcc"}
	cmd = append(cmd, args...)
	cmd = append(cmd, "-o", out)
	cmd = append(cmd, objects...)
	utils.Command(cmd, "Compiling")
}
