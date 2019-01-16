package index

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

/*
func getRootDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	rootDir := filepath.Join(cwd, "../..")
	return rootDir, nil
}
*/

func TestZip(t *testing.T) {
	t.Run("Zip()", func(t *testing.T) {
		rootDir, err := getRootDir()
		if err != nil {
			t.Error(err)
		}

		tmpPath := filepath.Join(rootDir, ".tmp")
		zipPath := filepath.Join(tmpPath, "sip.zip")

		err = Zip(rootDir, zipPath)
		if err != nil {
			t.Error(err)
		}

		os.RemoveAll(tmpPath)
	})
}

func TestUnzip(t *testing.T) {
	t.Run("Unzip()", func(t *testing.T) {
		rootDir, err := getRootDir()
		if err != nil {
			t.Error(err)
		}

		tmpPath := filepath.Join(rootDir, ".tmp")
		zipPath := filepath.Join(tmpPath, "sip.zip")
		extrctPath := filepath.Join(tmpPath, "unzip")

		err = Zip(rootDir, zipPath)
		if err != nil {
			t.Error(err)
		}

		err = Unzip(zipPath, extrctPath)
		if err != nil {
			t.Error(err)
		}

		os.RemoveAll(tmpPath)
	})
}

func TestGetGo(t *testing.T) {
	t.Run("getGo()", func(t *testing.T) {
		err := getGo()
		if err != nil {
			t.Error(err)
		}

		// Run twice to ensure that caching is considered
		err = getGo()
		if err != nil {
			t.Error(err)
		}
	})
}

func TestBuild(t *testing.T) {
	t.Run("Build()", func(t *testing.T) {
		os := runtime.GOOS
		arch := runtime.GOARCH

		goLoc = false

		err := Build(os, arch)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestUserAgent(t *testing.T) {
	t.Run("userAgent()", func(t *testing.T) {
		targetOS, targetArch := UserAgent("Browser/Version (" + runtime.GOOS + "; " + runtime.GOARCH + ";)")
		if targetOS != runtime.GOOS {
			t.Error("Failed to detect correct operating system.")
		}

		if targetArch != runtime.GOARCH {
			t.Error("Failed to detect correct architecture.")
		}
	})
}
