package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Println(("Usage:"))
	fmt.Println("  add <title> <description>")
	fmt.Println("  update <id> <title> <description>")
	fmt.Println("  delete <id>")
	fmt.Println("  start <id>")
	fmt.Println("  done <id>")
	fmt.Println("  list all|todo|inprogress|done")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "add":
	case "update":
	case "delete":
	case "list":
	default:
		usage()
	}
}
