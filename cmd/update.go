package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update devtime to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Checking for updates...")

		latest, err := fetchLatestVersion()
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		current := strings.TrimPrefix(Version, "v")
		latestClean := strings.TrimPrefix(latest, "v")

		if current == latestClean {
			fmt.Printf("Already up to date (v%s)\n", current)
			return nil
		}

		fmt.Printf("New version available: v%s → v%s\n", current, latestClean)

		archiveURL := buildDownloadURL(latest)
		fmt.Printf("Downloading %s...\n", archiveURL)

		archivePath, err := downloadToTemp(archiveURL)
		if err != nil {
			return fmt.Errorf("failed to download update: %w", err)
		}
		defer os.Remove(archivePath)

		binaryPath, err := extractBinary(archivePath)
		if err != nil {
			return fmt.Errorf("failed to extract update: %w", err)
		}
		defer os.Remove(binaryPath)

		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to locate current binary: %w", err)
		}
		execPath, err = filepath.EvalSymlinks(execPath)
		if err != nil {
			return fmt.Errorf("failed to resolve binary path: %w", err)
		}

		if err := replaceBinary(execPath, binaryPath); err != nil {
			return err
		}

		fmt.Printf("Successfully updated to v%s\n", latestClean)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

type githubRelease struct {
	TagName string `json:"tag_name"`
}

func fetchLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/arnaudhrt/devtime/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	if release.TagName == "" {
		return "", fmt.Errorf("no release tag found")
	}

	return release.TagName, nil
}

func buildDownloadURL(tag string) string {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	ext := ".tar.gz"
	if osName == "windows" {
		ext = ".zip"
	}
	return fmt.Sprintf(
		"https://github.com/arnaudhrt/devtime/releases/download/%s/devtime_%s_%s%s",
		tag, osName, arch, ext,
	)
}

func downloadToTemp(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	tmp, err := os.CreateTemp("", "devtime-update-*")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	return tmp.Name(), nil
}

func extractBinary(archivePath string) (string, error) {
	if runtime.GOOS == "windows" {
		return extractFromZip(archivePath)
	}
	return extractFromTarGz(archivePath)
}

func extractFromTarGz(archivePath string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if filepath.Base(hdr.Name) == "devtime" && hdr.Typeflag == tar.TypeReg {
			tmp, err := os.CreateTemp("", "devtime-bin-*")
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(tmp, tr); err != nil {
				tmp.Close()
				os.Remove(tmp.Name())
				return "", err
			}
			tmp.Close()
			if err := os.Chmod(tmp.Name(), 0755); err != nil {
				os.Remove(tmp.Name())
				return "", err
			}
			return tmp.Name(), nil
		}
	}

	return "", fmt.Errorf("devtime binary not found in archive")
}

func extractFromZip(archivePath string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, f := range r.File {
		if filepath.Base(f.Name) == "devtime.exe" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			tmp, err := os.CreateTemp("", "devtime-bin-*.exe")
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(tmp, rc); err != nil {
				tmp.Close()
				os.Remove(tmp.Name())
				return "", err
			}
			tmp.Close()
			return tmp.Name(), nil
		}
	}

	return "", fmt.Errorf("devtime.exe not found in archive")
}

func replaceBinary(dst, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	info, err := os.Stat(dst)
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied updating %s — try running with sudo", dst)
		}
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to write updated binary: %w", err)
	}

	return nil
}
