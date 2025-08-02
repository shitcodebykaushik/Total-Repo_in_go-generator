package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

const GITHUB_API = "https://api.github.com/search/repositories?q=language:go&per_page=100"

type Repository struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    URL         string `json:"html_url"`
}

type GitHubResponse struct {
    TotalCount       int          `json:"total_count"`
    IncompleteResult bool         `json:"incomplete_results"`
    Items            []Repository `json:"items"`
}

func main() {
    file, err := os.Create("README.md")
    if err != nil {
        log.Fatalf("Error encountered while creating file: %v", err)
    }
    defer file.Close()

    fmt.Fprintf(file, "# Tools created in Go :gopher: \n\n")
    fmt.Fprintf(file, "This README is a dynamic list of popular Go projects fetched from the GitHub API.\n\n")
    fmt.Fprintf(file, "Inspired by the 'rusty-alternatives' project.\n\n")

    resp, err := http.Get(GITHUB_API)
    if err != nil {
        log.Fatalf("Error making HTTP request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Fatalf("GitHub API returned non-200 status: %s", resp.Status)
    }

    var githubResponse GitHubResponse
    err = json.NewDecoder(resp.Body).Decode(&githubResponse)
    if err != nil {
        log.Fatalf("Error decoding JSON response: %v", err)
    }

    fmt.Fprintf(file, "Total repositories created in Go: %d\n\n", githubResponse.TotalCount)

    for _, item := range githubResponse.Items {
        description := item.Description
        if description == "" {
            description = "No description provided."
        }
        fmt.Fprintf(file, "- [%s](%s) - %s\n", item.Name, item.URL, description)
    }

    fmt.Fprintf(file, "\n\n## License\n\n")
    fmt.Fprintf(file, "This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details\n\n")
    fmt.Fprintf(file, "## Acknowledgments\n\n")
    fmt.Fprintf(file, "- [GitHub](https://github.com)\n")
    fmt.Fprintf(file, "- [Go Language](https://golang.org)\n\n")

    now := time.Now().Format("02-01-2006 15:04:05")
    fmt.Fprintf(file, "##### _Last Run on %s_", now)

    log.Println("README.md file generated successfully.")
}