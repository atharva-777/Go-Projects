package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/atharva-777/go-projects/task-tracker/store"
)

func usage() {
	fmt.Println(("Usage:"))
	fmt.Println("  add <title> <description>")
	fmt.Println("  update <id> <title> <description>")
	fmt.Println("  delete <id>")
	fmt.Println("  setStatus <id>")
	fmt.Println("  list all|todo|inprogress|done")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	cmd := os.Args[1]
	db := store.New("tasks.json")

	switch cmd {
	case "add":
		if len(os.Args) < 4 {
			usage()
			return
		}
		title := os.Args[2]
		desc := os.Args[3]
		t := db.Add(title, desc)
		fmt.Printf("Added: %d %s\n", t.ID, t.Title)
	case "update":
		if len(os.Args) < 5 {
			usage()
			return
		}
		id, _ := strconv.Atoi(os.Args[2])
		title := os.Args[3]
		desc := os.Args[4]
		if db.Update(id, title, desc) {
			fmt.Printf("Updated")
		} else {
			fmt.Printf("Not found")
		}
	case "delete":
		id, _ := strconv.Atoi(os.Args[2])

		if db.Delete(id) {
			fmt.Println("Deleted")
		} else {
			fmt.Println("Not found")
		}
	case "setStatus":
		id, _ := strconv.Atoi(os.Args[2])
		status := os.Args[3]
		if db.SetStatus(id, status) {
			fmt.Println("Changed status")
		} else {
			fmt.Println("Not found")
		}
	case "list":
		filter := "all"
		if len(os.Args) >= 3 {
			filter = os.Args[2]
		}
		for _, t := range db.List(filter) {
			fmt.Printf("%d [%s] %s - %s\n", t.ID, t.Status, t.Title, t.Description)
		}
	default:
		usage()
	}
}
