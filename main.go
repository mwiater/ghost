package main

import (
	"github.com/mwiater/ghost/cmd"
	"github.com/mwiater/ghost/utils"
)

func main() {
	utils.ClearTerminal()

	cmd.Execute()
}
