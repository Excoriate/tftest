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

// FetchReleases fetches GitHub releases for a given repo URL.
// If `onlyLatest` is true, it returns the most recent release.
// Otherwise, it returns all releases.
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

	releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error fetching releases: %w", err)
	}

	return releases, nil
}

// GetOwnerAndRepoFromURL extracts the owner and repo name from a GitHub URL
func GetOwnerAndRepoFromURL(repoURL string) (string, string, error) {
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
