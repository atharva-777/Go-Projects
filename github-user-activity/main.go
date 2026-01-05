package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide username!")
		return
	}

	username := os.Args[1]
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build request", err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "github-activity-cli")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Fprintln(os.Stderr, "request error ", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		// ok
	case 404:
		fmt.Println("user not found")
		return
	case 403:
		fmt.Println("forbidden or rate limited ", resp.Status)
		return
	default:
		fmt.Println("http error:", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Fprintln(os.Stderr, "read body error:", err)
		os.Exit(1)
	}

	var pretty bytes.Buffer

	if err := json.Indent(&pretty, body, "", " "); err != nil {
		fmt.Println(string(body))
		return
	}

	fmt.Println(pretty.String())
}
