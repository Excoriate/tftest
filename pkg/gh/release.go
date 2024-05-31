package gh

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

// FetchReleases fetches GitHub releases for a given repository URL.
// If `onlyLatest` is true, it returns the most recent release. Otherwise, it returns all releases.
//
// Parameters:
//   - repoURL: The URL of the GitHub repository (e.g., "https://github.com/owner/repo").
//   - onlyLatest: A boolean flag indicating whether to fetch only the latest release.
//
// Returns:
//   - []*github.RepositoryRelease: A slice of GitHub repository releases.
//   - error: An error if the releases could not be fetched.
//
// Example:
//
//	releases, err := FetchReleases("https://github.com/owner/repo", true)
//	if err != nil {
//	    log.Fatalf("Error fetching releases: %v", err)
//	}
//	for _, release := range releases {
//	    fmt.Printf("Release: %s\n", *release.TagName)
//	}
func FetchReleases(repoURL string, onlyLatest bool) ([]*github.RepositoryRelease, error) {
	ctx := context.Background()

	// Setup GitHub client with authentication
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	))
	client := github.NewClient(tc)

	owner, repo, err := GetOwnerAndRepoFromURL(repoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GitHub URL: %w", err)
	}

	if onlyLatest {
		release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
		if err != nil {
			return nil, fmt.Errorf("error fetching the latest release: %w", err)
		}
		return []*github.RepositoryRelease{release}, nil
	}

	releases, _, releaseErr := client.Repositories.ListReleases(ctx, owner, repo, &github.ListOptions{})
	if releaseErr != nil {
		return nil, fmt.Errorf("error fetching releases: %w", releaseErr)
	}

	return releases, nil
}

// GetOwnerAndRepoFromURL extracts the owner and repository name from a GitHub URL.
//
// Parameters:
//   - repoURL: The URL of the GitHub repository (e.g., "https://github.com/owner/repo").
//
// Returns:
//   - owner: The owner of the repository.
//   - repo: The name of the repository.
//   - err: An error if the URL could not be parsed or is in an invalid format.
//
// Example:
//
//	owner, repo, err := GetOwnerAndRepoFromURL("https://github.com/owner/repo")
//	if err != nil {
//	    log.Fatalf("Error parsing URL: %v", err)
//	}
//	fmt.Printf("Owner: %s, Repo: %s\n", owner, repo)
func GetOwnerAndRepoFromURL(repoURL string) (owner, repo string, err error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", fmt.Errorf("error parsing URL: %w", err)
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid GitHub repo URL format")
	}
	return parts[0], parts[1], nil
}
