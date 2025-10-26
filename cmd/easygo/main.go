package main

import (
	"easygo/internal/cli"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		// CLI mode
		cli.Execute()
	} else {
		// Web mode - show usage
		fmt.Println("EasyGo Panel v1.0.0")
		fmt.Println("Usage:")
		fmt.Println("  easygo web          Start web panel")
		fmt.Println("  easygo <command>    Run CLI command")
		fmt.Println("  easygo help         Show available commands")
	}
}