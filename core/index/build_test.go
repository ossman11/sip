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

func TestTar(t *testing.T) {
	t.Run("Tar()", func(t *testing.T) {
		rootDir, err := getRootDir()
		if err != nil {
			t.Error(err)
		}

		tmpPath := filepath.Join(rootDir, ".tmp")
		zipPath := filepath.Join(tmpPath, "sip.tar.gz")

		err = Tar(rootDir, zipPath)
		if err != nil {
			t.Error(err)
		}

		os.RemoveAll(tmpPath)
	})
}

func TestUntar(t *testing.T) {
	t.Run("Untar()", func(t *testing.T) {
		rootDir, err := getRootDir()
		if err != nil {
			t.Error(err)
		}

		tmpPath := filepath.Join(rootDir, ".tmp")
		zipPath := filepath.Join(tmpPath, "sip.tar.gz")
		extrctPath := filepath.Join(tmpPath, "unzip")

		err = Tar(rootDir, zipPath)
		if err != nil {
			t.Error(err)
		}

		err = Untar(zipPath, extrctPath)
		if err != nil {
			t.Error(err)
		}

		os.RemoveAll(tmpPath)
	})
}

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

		rootDir, err := getRootDir()
		if err != nil {
			t.Error(err)
		}

		tmpPath := filepath.Join(rootDir, ".tmp")
		os.RemoveAll(tmpPath)
	})
}

func TestBuild(t *testing.T) {
	t.Run("Build()", func(t *testing.T) {
		goLoc = false

		err := Build(runtime.GOOS, runtime.GOARCH)
		if err != nil {
			t.Error(err)
		}

		rootDir, err := getRootDir()
		if err != nil {
			t.Error(err)
		}

		tmpPath := filepath.Join(rootDir, ".tmp")
		os.RemoveAll(tmpPath)
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
