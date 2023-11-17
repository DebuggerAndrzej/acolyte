package backend

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func RunCommand(command string) {
	cmd := exec.Command("fish", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()

	buf := bufio.NewReader(stdout)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			os.Exit(0)
		}
		fmt.Println(string(line))
	}
}
