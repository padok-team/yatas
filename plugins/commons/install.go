package commons

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/go-github/v35/github"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2"
)

func (c *Plugin) Install() (string, error) {

	dir, err := homedir.Expand("~/.yatas.d/plugins")
	if err != nil {
		return "", fmt.Errorf("Failed to get plugin dir: %w", err)
	}

	path := filepath.Join(dir, c.InstallPath()+fileExt())
	log.Printf("[DEBUG] Mkdir plugin dir: %s", filepath.Dir(path))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("Failed to mkdir to %s: %w", filepath.Dir(path), err)
	}

	assets, err := c.fetchReleaseAssets()
	if err != nil {
		return "", fmt.Errorf("Failed to fetch GitHub releases: %w", err)
	}

	log.Printf("[DEBUG] Download checksums.txt")
	checksumsFile, err := c.downloadToTempFile(assets["checksums.txt"])
	fmt.Println("Debug 1")
	if checksumsFile != nil {
		defer os.Remove(checksumsFile.Name())
	}
	if err != nil {
		return "", fmt.Errorf("Failed to download checksums.txt: %s", err)
	}

	log.Printf("[DEBUG] Download %s", c.AssetName())
	zipFile, err := c.downloadToTempFile(assets[c.AssetName()])
	if zipFile != nil {
		defer os.Remove(zipFile.Name())
	}
	if err != nil {
		return "", fmt.Errorf("Failed to download %s: %s", c.AssetName(), err)
	}

	if err = extractFileFromZipFile(zipFile, path); err != nil {
		return "", fmt.Errorf("Failed to extract binary from %s: %s", c.AssetName(), err)
	}

	log.Printf("[DEBUG] Installed %s successfully", path)
	return path, nil
}

func GetRelease(ctx context.Context, client *github.Client, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error) {
	if tag == "latest" {
		return client.Repositories.GetLatestRelease(ctx, owner, repo)
	} else {
		return client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	}
}

func GetLatestReleaseTag(plugin Plugin) (string, error) {
	ctx := context.Background()
	client := newGitHubClient(ctx)
	plugin.Validate()
	latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, plugin.SourceOwner, plugin.SourceRepo)
	if err != nil {
		return "", err
	}

	return latestRelease.GetName()[1:], nil
}

// fetchReleaseAssets fetches assets from the GitHub release.
// The release is determined by the source path and tag name.
func (c *Plugin) fetchReleaseAssets() (map[string]*github.ReleaseAsset, error) {
	assets := map[string]*github.ReleaseAsset{}

	ctx := context.Background()
	client := newGitHubClient(ctx)

	log.Printf("[DEBUG] Request to https://api.github.com/repos/%s/%s/releases/tags/%s", c.SourceOwner, c.SourceRepo, c.TagName())
	release, _, err := GetRelease(ctx, client, c.SourceOwner, c.SourceRepo, c.TagName())
	if err != nil {
		return assets, err
	}

	for _, asset := range release.Assets {
		log.Printf("[DEBUG] asset found: %s", asset.GetName())
		assets[asset.GetName()] = asset

	}
	return assets, nil
}

// downloadToTempFile download assets from GitHub to a local temp file.
// It is the caller's responsibility to delete the generated the temp file.
func (c *Plugin) downloadToTempFile(asset *github.ReleaseAsset) (*os.File, error) {
	if asset == nil {
		return nil, fmt.Errorf("file not found in the GitHub release. Does the release contain the file with the correct name ?")
	}

	ctx := context.Background()
	client := newGitHubClient(ctx)

	log.Printf("[DEBUG] Request to https://api.github.com/repos/%s/%s/releases/assets/%d", c.SourceOwner, c.SourceRepo, asset.GetID())
	downloader, _, err := client.Repositories.DownloadReleaseAsset(ctx, c.SourceOwner, c.SourceRepo, asset.GetID(), http.DefaultClient)
	if err != nil {
		return nil, err
	}
	fmt.Println("Debug 2")
	file, err := os.CreateTemp("", "yatas-tmp-*")
	if err != nil {
		return nil, err
	}
	fmt.Println("Debug 3")

	fmt.Println(file.Name())
	fmt.Println(downloader)
	if _, err = io.Copy(file, downloader); err != nil {
		fmt.Println("Debug 4")
		return file, err
	}
	downloader.Close()
	if _, err := file.Seek(0, 0); err != nil {
		return file, err
	}

	log.Printf("[DEBUG] Downloaded to %s", file.Name())
	return file, nil
}

func extractFileFromZipFile(zipFile *os.File, savePath string) error {
	zipFileStat, err := zipFile.Stat()
	if err != nil {
		return err
	}
	zipReader, err := zip.NewReader(zipFile, zipFileStat.Size())
	if err != nil {
		return err
	}

	var reader io.ReadCloser
	for _, f := range zipReader.File {
		log.Printf("[DEBUG] file found in zip: %s", f.Name)
		if f.Name != filepath.Base(savePath) {
			continue
		}

		reader, err = f.Open()
		if err != nil {
			return err
		}
		break
	}
	if reader == nil {
		return fmt.Errorf("file not found. Does the zip contain %s ?", filepath.Base(savePath))
	}

	file, err := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		os.Remove(file.Name())
		return err
	}

	return nil
}

func newGitHubClient(ctx context.Context) *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return github.NewClient(nil)
	}

	// log.Printf("[DEBUG] GITHUB_TOKEN set, plugin requests to the GitHub API will be authenticated")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	return github.NewClient(oauth2.NewClient(ctx, ts))
}

func fileExt() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}
