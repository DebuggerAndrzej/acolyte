package main

import (
	"github.com/DebuggerAndrzej/acolyte/backend"
)

func main() {
	backend.RunCommand("ping -c 3 google.pl")
}
