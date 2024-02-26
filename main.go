package main

import (
    "context"
    "encoding/xml"
    "fmt"
    "log"
    "net/http"
    "os"
//GitHub Oth dependenies 
    "github.com/google/go-github/v35/github"
    "golang.org/x/oauth2"
)

func main() {
    // GitHub personal access token
    githubToken := "5e4abf421189582dc5749d0b1bf6dafb23709179"

    // Create a GitHub client with OAuth2 authentication
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: githubToken},
    )
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)

    // List all repositories for the authenticated user
    repos, _, err := client.Repositories.List(ctx, "", nil)
    if err != nil {
        log.Fatal(err)
    }

    // Select a repository from the list
    selectedRepo := repos[0] // Choose the first repository for demonstration

    // Fetch the contents of the repository
    contents, _, err := client.Repositories.GetContents(ctx, selectedRepo.GetOwner().GetLogin(), selectedRepo.GetName(), "", nil)
    if err != nil {
        log.Fatal(err)
    }

    // Iterate through the contents to find and parse pom.xml files
    for _, content := range contents {
        if *content.Type == "file" && *content.Name == "pom.xml" {
            // Fetch the file contents
            fileContent, _, _, err := client.Repositories.GetContents(ctx, selectedRepo.GetOwner().GetLogin(), selectedRepo.GetName(), *content.Path, nil)
            if err != nil {
                log.Fatal(err)
            }

            // Parse the XML content to extract dependencies
            dependencies, err := parsePomXML(fileContent.GetContent())
            if err != nil {
                log.Fatal(err)
            }

            // Print dependencies and versions
            for _, dependency := range dependencies {
                fmt.Printf("%s: Version %s\n", dependency.ArtifactID, dependency.Version)
            }
        }
    }
}

// Dependency represents a Maven dependency
type Dependency struct {
    GroupID    string `xml:"groupId"`
    ArtifactID string `xml:"artifactId"`
    Version    string `xml:"version"`
}

// parsePomXML parses the pom.xml file and extracts dependencies
func parsePomXML(xmlContent string) ([]Dependency, error) {
    // Define a struct to hold the parsed XML data
    type Project struct {
        Dependencies []Dependency `xml:"dependencies>dependency"`
    }

    // Parse the XML content
    var project Project
    if err := xml.Unmarshal([]byte(xmlContent), &project); err != nil {
        return nil, err
    }

    return project.Dependencies, nil
}