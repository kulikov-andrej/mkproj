package installer

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var runtimeGOOS = runtime.GOOS
var runtimeGOARCH = runtime.GOARCH

const (
	repository = "kulikov-andrej/mkproj"
	apiURL     = "https://api.github.com/repos/" + repository + "/releases/latest"
)

type release struct {
	TagName string  `json:"tag_name"`
	Assets  []asset `json:"assets"`
}

type asset struct {
	Name   string `json:"name"`
	URL    string `json:"browser_download_url"`
	Digest string `json:"digest"`
}

func Run(output io.Writer) error {
	if runtimeGOOS != "windows" &&
		runtimeGOOS != "linux" &&
		runtimeGOOS != "darwin" {
		return fmt.Errorf("unsupported platform: %s", runtimeGOOS)
	}

	fmt.Fprintln(output, "Installing mkproj...")

	release, err := fetchRelease()
	if err != nil {
		return err
	}

	asset, err := findAsset(release)
	if err != nil {
		return err
	}

	tempFile, err := os.CreateTemp("", "mkproj-download-*")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}

	fmt.Fprintf(output, "Downloading %s...\n", release.TagName)
	if err := download(asset.URL, tempPath); err != nil {
		return err
	}
	if err := verify(tempPath, asset.Digest); err != nil {
		return err
	}

	paths, err := install(tempPath)
	if err != nil {
		return err
	}

	fmt.Fprintln(output)
	fmt.Fprintln(output, "Installed:")
	fmt.Fprintf(output, "  %s\n", paths.Binary)
	if paths.CreatedTemplates {
		fmt.Fprintln(output)
		fmt.Fprintln(output, "Created template directory:")
		fmt.Fprintf(output, "  %s\n", paths.Templates)
	}

	if err := reportVersion(paths.Binary, output); err != nil {
		return err
	}

	return reportPath(paths, output)
}

func fetchRelease() (release, error) {
	request, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return release{}, err
	}
	request.Header.Set("Accept", "application/vnd.github+json")

	response, err := (&http.Client{}).Do(request)
	if err != nil {
		return release{}, fmt.Errorf("fetch latest release: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return release{}, fmt.Errorf("fetch latest release: %s", response.Status)
	}

	var result release
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return release{}, fmt.Errorf("decode release metadata: %w", err)
	}
	return result, nil
}

func findAsset(release release) (asset, error) {
	name := fmt.Sprintf(
		"mkproj-%s-%s",
		runtimeGOOS,
		runtimeGOARCH,
	)
	if runtimeGOOS == "windows" {
		name += ".exe"
	}

	for _, candidate := range release.Assets {
		if candidate.Name == name {
			if !strings.HasPrefix(candidate.Digest, "sha256:") {
				return asset{}, fmt.Errorf("release asset has no supported SHA-256 digest")
			}
			return candidate, nil
		}
	}

	return asset{}, fmt.Errorf("release asset not found: %s", name)
}

func download(url, path string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download release asset: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("download release asset: %s", response.Status)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create download: %w", err)
	}
	defer file.Close()
	if _, err := io.Copy(file, response.Body); err != nil {
		return fmt.Errorf("write download: %w", err)
	}
	return nil
}

func verify(path, digest string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open download: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("hash download: %w", err)
	}
	actual := hex.EncodeToString(hash.Sum(nil))
	expected := strings.TrimPrefix(digest, "sha256:")
	if !strings.EqualFold(actual, expected) {
		return fmt.Errorf("SHA-256 verification failed")
	}
	return nil
}

type installPaths struct {
	Binary           string
	Templates        string
	CreatedTemplates bool
}

func copyFile(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer output.Close()
	_, err = io.Copy(output, input)
	return err
}

func templateDir(configDir string) string {
	return filepath.Join(configDir, "mkproj", "templates")
}

func reportVersion(path string, output io.Writer) error {
	command := exec.Command(path, "version")
	command.Stdout = output
	command.Stderr = output
	if err := command.Run(); err != nil {
		return fmt.Errorf("run installed binary: %w", err)
	}
	return nil
}
