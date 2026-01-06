package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

const userAgent = "github-activity-cli"

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  github-activity <username>           # show recent events (default)")
	fmt.Println("  github-activity repos <username>     # list public repos")
	fmt.Println("  github-activity profile <username>   # show user profile")
	fmt.Println("  github-activity repo-langs <owner> <repo>  # show languages for a repo")
}

func fetch(url string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	// Use only unauthenticated public API (no token)
	req.Header.Set("User-Agent", userAgent)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

type Profile struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	Bio         string `json:"bio"`
	Location    string `json:"location"`
	HTMLURL     string `json:"html_url"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	CreatedAt   string `json:"created_at"`
}

type Repo struct {
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	HTMLURL         string `json:"html_url"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	Language        string `json:"language"`
	UpdatedAt       string `json:"updated_at"`
}

func handleProfile(username string) error {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	body, status, err := fetch(url)
	if err != nil {
		return err
	}
	switch status {
	case 200:
	case 404:
		fmt.Println("user not found")
		return nil
	case 403:
		fmt.Println("forbidden or rate limited:", status)
		return nil
	default:
		fmt.Println("http error:", status)
		return nil
	}
	var p Profile
	if err := json.Unmarshal(body, &p); err != nil {
		return err
	}
	fmt.Printf("%s (%s)\n", p.Name, p.Login)
	fmt.Printf("Bio: %s\n", p.Bio)
	fmt.Printf("Location: %s\n", p.Location)
	fmt.Printf("Repos: %d  Followers: %d  Following: %d\n", p.PublicRepos, p.Followers, p.Following)
	fmt.Printf("Profile: %s\n", p.HTMLURL)
	return nil
}

func handleRepos(username string) error {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100", username)
	body, status, err := fetch(url)
	if err != nil {
		return err
	}
	switch status {
	case 200:
	case 404:
		fmt.Println("user not found")
		return nil
	case 403:
		fmt.Println("forbidden or rate limited:", status)
		return nil
	default:
		fmt.Println("http error:", status)
		return nil
	}
	var repos []Repo
	if err := json.Unmarshal(body, &repos); err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTARS\tFORKS\tLANG\tUPDATED")
	for _, r := range repos {
		fmt.Fprintf(w, "%s\t%d\t%d\t%s\t%s\n", r.FullName, r.StargazersCount, r.ForksCount, r.Language, r.UpdatedAt)
	}
	return w.Flush()
}

func handleRepoLangs(owner, repo string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", owner, repo)
	body, status, err := fetch(url)
	if err != nil {
		return err
	}
	switch status {
	case 200:
	case 404:
		fmt.Println("repo not found")
		return nil
	case 403:
		fmt.Println("forbidden or rate limited:", status)
		return nil
	default:
		fmt.Println("http error:", status)
		return nil
	}
	var langs map[string]int
	if err := json.Unmarshal(body, &langs); err != nil {
		return err
	}
	if len(langs) == 0 {
		fmt.Println("no languages data")
		return nil
	}
	for k, v := range langs {
		fmt.Printf("%s: %d\n", k, v)
	}
	return nil
}

func showEvents(username string) error {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	body, status, err := fetch(url)
	if err != nil {
		return err
	}
	switch status {
	case 200:
	case 404:
		fmt.Println("user not found")
		return nil
	case 403:
		fmt.Println("forbidden or rate limited:", status)
		return nil
	default:
		fmt.Println("http error:", status)
		return nil
	}
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, body, "", "  "); err != nil {
		fmt.Println(string(body))
		return nil
	}
	fmt.Println(pretty.String())
	return nil
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "repos":
		if len(os.Args) < 3 {
			fmt.Println("repos requires username")
			return
		}
		_ = handleRepos(os.Args[2])
	case "profile":
		if len(os.Args) < 3 {
			fmt.Println("profile requires username")
			return
		}
		_ = handleProfile(os.Args[2])
	case "repo-langs":
		if len(os.Args) < 4 {
			fmt.Println("repo-langs requires owner and repo")
			return
		}
		_ = handleRepoLangs(os.Args[2], os.Args[3])
	case "-h", "--help", "help":
		usage()
	default:
		// treat as username for events
		_ = showEvents(cmd)
	}
}
