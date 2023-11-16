package main

import (
	"fmt"
	"os/exec"
)

func main() {
	command := "ls -la"
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))
}
