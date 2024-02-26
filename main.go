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
    // Provide your GitHub personal access token
    githubToken := "YOUR_GITHUB_TOKEN"

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

/*Here's a step-by-step guide to set up Go and complete the project:

### Step 1: Install Go

1. Download and install Go from the official website: [https://golang.org/dl/](https://golang.org/dl/)
2. Follow the installation instructions provided for your operating system.

### Step 2: Set up a GitHub Personal Access Token

1. Go to your GitHub account settings.
2. Navigate to Developer settings > Personal access tokens.
3. Click on "Generate new token" and provide necessary permissions (e.g., `repo` scope).
4. Copy the generated token.

### Step 3: Set up Go Modules (Optional but recommended)

1. Go modules enable dependency management within your Go projects.
2. Enable Go modules by setting the environment variable `GO111MODULE` to `on`:
    ```bash
    export GO111MODULE=on
    ```

### Step 4: Install Required Packages

1. Open a terminal or command prompt.
2. Install the required Go packages using the `go get` command:
    ```bash
    go get github.com/google/go-github/v35/github
    go get golang.org/x/oauth2
    ```

### Step 5: Write the Go Code

1. Open your favorite text editor or integrated development environment (IDE).
2. Write the Go code provided in the previous response. Save the file with a `.go` extension, for example, `main.go`.

### Step 6: Update the Code with Your GitHub Token

1. Replace `"YOUR_GITHUB_TOKEN"` in the code with the personal access token you generated earlier.

### Step 7: Implement XML Parsing Logic

1. Implement the `parsePomXML` function to parse the `pom.xml` files and extract dependencies. You can use Go's `encoding/xml` package or other XML parsing libraries.

### Step 8: Test and Run the Code

1. Navigate to the directory containing your Go code (`main.go`).
2. Run the code using the `go run` command:
    ```bash
    go run main.go
    ```

### Step 9: Troubleshooting and Debugging

1. If you encounter any errors or issues, carefully review the error messages and consult the documentation for the GitHub API and the Go libraries being used.
2. Make sure your GitHub personal access token has the necessary permissions.
3. Ensure that your XML parsing logic is correctly implemented.

### Step 10: Iterate and Improve

1. Test the parser with different repositories and handle edge cases.
2. Refactor and improve your code as necessary based on feedback and additional requirements.

By following these steps, you should be able to set up Go, authenticate with the GitHub API, retrieve repository information, fetch `pom.xml` files, parse XML content, and extract dependencies as per your project requirements.*/