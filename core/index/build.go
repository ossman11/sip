package index

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Unzip(src, dest string) error {
	// Open zip archive
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create root destination directory
	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
		} else {
			os.MkdirAll(filepath.Dir(path), 0755)
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func Untar(src, dest string) error {
	tarFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	gzr, err := gzip.NewReader(tarFile)
	if err != nil {
		return err
	}
	defer gzr.Close()

	// Open tar archive
	r := tar.NewReader(gzr)

	// Create root destination directory
	os.MkdirAll(dest, 0755)

	// Closure to keep unzip and untar simulair
	extractAndWriteFile := func(h *tar.Header) error {
		path := filepath.Join(dest, h.Name)

		if h.FileInfo().IsDir() {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
		} else {
			err := os.MkdirAll(filepath.Dir(path), 0755)
			if err != nil {
				return err
			}

			f, err := os.Create(path)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, r)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for {
		hdr, err := r.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = extractAndWriteFile(hdr)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadIgnore(src string) []string {
	file, err := ioutil.ReadFile(src)
	if err != nil {
		return nil
	}
	ignoreContent := string(file)
	rules := strings.Split(ignoreContent, "\n")
	ret := []string{}
	for _, rule := range rules {
		rule = strings.Trim(rule, " ")
		if strings.Index(rule, "#") == 0 || rule == "" {
			continue
		}

		ret = append(ret, rule)
	}

	return ret
}

func compressAndReadFileZip(src, pth string, filters []string, w *zip.Writer, f os.FileInfo) error {
	fPth := filepath.Join(pth, f.Name())

	for _, filter := range filters {
		m, _ := filepath.Match(filter, f.Name())
		if !m {
			rPth, err := filepath.Rel(src, fPth)
			if err == nil {
				m, _ = filepath.Match(filter, rPth)
			}
		}
		if m {
			return nil
		}
	}

	if !f.Mode().IsRegular() {
		newFilter := loadIgnore(filepath.Join(fPth, ".gitignore"))
		filters = append(filters, newFilter...)

		files, err := ioutil.ReadDir(fPth)
		if err != nil {
			return err
		}

		for _, file := range files {
			err := compressAndReadFileZip(src, fPth, filters, w, file)
			if err != nil {
				return err
			}
		}
	} else {
		relName, err := filepath.Rel(src, fPth)
		if err != nil {
			return err
		}

		fw, err := w.Create(relName)
		if err != nil {
			return err
		}

		fr, err := os.Open(fPth)
		if err != nil {
			return err
		}
		defer fr.Close()

		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}
	}
	return nil
}

func Zip(src, dest string) error {
	// Create a buffer to write our archive to.
	os.MkdirAll(filepath.Dir(dest), 0755)
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create a new zip archive.
	w := zip.NewWriter(out)
	defer w.Close()

	file, err := os.Stat(src)
	if err != nil {
		return err
	}

	par := filepath.Join(src, "..")

	return compressAndReadFileZip(par, par, []string{".git"}, w, file)
}

func compressAndReadFileTar(src, pth string, filters []string, w *tar.Writer, f os.FileInfo) error {
	fPth := filepath.Join(pth, f.Name())

	for _, filter := range filters {
		m, _ := filepath.Match(filter, f.Name())
		if !m {
			rPth, err := filepath.Rel(src, fPth)
			if err == nil {
				m, _ = filepath.Match(filter, rPth)
			}
		}
		if m {
			return nil
		}
	}

	relName, err := filepath.Rel(src, fPth)

	// create a new dir/file header
	header, err := tar.FileInfoHeader(f, relName)
	if err != nil {
		return err
	}

	header.Name = strings.TrimPrefix(relName, string(filepath.Separator))

	if err := w.WriteHeader(header); err != nil {
		return err
	}

	if !f.Mode().IsRegular() {
		newFilter := loadIgnore(filepath.Join(fPth, ".gitignore"))
		filters = append(filters, newFilter...)

		files, err := ioutil.ReadDir(fPth)
		if err != nil {
			return err
		}

		for _, file := range files {
			err := compressAndReadFileTar(src, fPth, filters, w, file)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// open files for taring
	file, err := os.Open(fPth)
	if err != nil {
		return err
	}
	defer file.Close()

	// copy file data into tar writer
	if _, err := io.Copy(w, file); err != nil {
		return err
	}

	return nil
}

func Tar(src, dest string) error {
	// ensure the src actually exists before trying to tar it
	file, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create a buffer to write our archive to.
	os.MkdirAll(filepath.Dir(dest), 0755)
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	gzw := gzip.NewWriter(out)
	defer gzw.Close()

	w := tar.NewWriter(gzw)
	defer w.Close()

	par := filepath.Join(src, "..")

	// walk path
	return compressAndReadFileTar(par, par, []string{".git"}, w, file)
}

var mappings = []struct {
	match string
	arch  string
	os    string
}{
	{os: "android"},
	{os: "darwin"},
	{os: "dragonfly"},
	{os: "freebsd"},
	{os: "linux"},
	{os: "netbsd"},
	{os: "openbsd"},
	{os: "plan9"},
	{os: "solaris"},
	{os: "windows"},
	{arch: "arm"},
	{arch: "386"},
	{arch: "amd64"},
	{arch: "ppc64"},
	{arch: "ppc64le"},
	{arch: "mips"},
	{arch: "mipsle"},
	{arch: "mips64"},
	{arch: "mips64le"},
	{
		match: "x86_64",
		arch:  "amd64",
	},
	{
		match: "x64",
		arch:  "amd64",
	},
	{
		match: "win64",
		arch:  "amd64",
	},
}

func UserAgent(userAgent string) (string, string) {
	root := true
	buff := ""
	var targetOS string
	var targetArch string

	for _, cv := range userAgent {
		if cv == ' ' && root || cv == ';' && !root {
			buff = strings.Trim(buff, " ;()")
			buff = strings.ToLower(buff)

			for _, cm := range mappings {
				if (cm.match != "" && strings.Index(buff, cm.match) > -1) ||
					(cm.arch != "" && strings.Index(buff, cm.arch) > -1) ||
					(cm.os != "" && strings.Index(buff, cm.os) > -1) {
					if cm.arch != "" {
						targetArch = cm.arch
					}
					if cm.os != "" {
						targetOS = cm.os
					}
				}
			}

			buff = ""
			continue
		}

		if cv == '(' && root || cv == ')' && !root {
			root = !root
			continue
		}

		buff = buff + string(cv)
	}

	return targetOS, targetArch
}

var locTmp = ".tmp"
var goLoc = false
var goEns = false
var goLocDir = filepath.Join(locTmp, "go")
var goLocSip = filepath.Join(goLocDir, "src", "github.com", "ossman11", "sip")
var goLocBin = filepath.Join(goLocDir, "bin")

func localize(pth string) string {
	root, err := getRootDir()
	if err != nil {
		return pth
	}
	res := filepath.Join(root, pth)
	return res
}

func getGo() error {
	isTar := runtime.GOOS != "windows"
	archiveExtension := ".tar.gz"
	if !isTar {
		archiveExtension = ".zip"
	}
	goLocArc := localize(goLocDir) + archiveExtension

	file, err := os.Open(localize(goLocBin))
	if err == nil {
		file.Close()
		return nil
	}

	err = os.MkdirAll(localize(locTmp), 0755)
	file, err = os.Create(goLocArc)
	if err != nil {
		return err
	}
	defer file.Close()

	dlURL := "https://dl.google.com/go/" +
		runtime.Version() + "." + runtime.GOOS + "-" + runtime.GOARCH + archiveExtension
	res, err := http.Get(dlURL)
	if err != nil {
		return err
	}

	io.Copy(file, res.Body)

	if isTar {
		err = Untar(goLocArc, localize(locTmp))
	} else {
		err = Unzip(goLocArc, localize(locTmp))
	}
	if err != nil {
		return err
	}

	return nil
}

func getGoCommand() string {
	goStr := "go"

	if goLoc {
		goStr = filepath.Join(localize(goLocBin), goStr)
	}

	return goStr
}

func ensureGo() {
	if goEns {
		return
	}

	cmd := exec.Command(getGoCommand(), "version")
	err := cmd.Run()
	if err != nil {
		goLoc = true
		err := getGo()
		if err != nil {
			log.Fatal(err)
		}
		ensureGo()
		return
	}
}

func runGo(args, env []string, stdout, stderr *bytes.Buffer) error {
	rootDir, err := getRootDir()
	if err != nil {
		rootDir = ""
	}

	if goLoc {
		env = append(env, "GOROOT="+localize(goLocDir))
	}

	env = append(os.Environ(), env...)

	ensureGo()
	cmd := exec.Command(getGoCommand(), args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = env
	cmd.Dir = rootDir

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func getRootDir() (string, error) {
	str, err := os.Getwd()
	if err != nil {
		return "", err
	}

	splt := strings.Split(str, "sip")

	if len(splt) < 2 {
		return str, nil
	}

	subDirs := strings.Split(splt[1], string(filepath.Separator))
	subCnt := len(subDirs)

	for subCnt > 1 {
		str = filepath.Join(str, "..")
		subCnt--
	}

	return str, nil
}

func Build(tos, tarch string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	pckStr := "github.com/ossman11/sip"

	root, err := getRootDir()
	if err != nil {
		return err
	}

	if goLoc {
		target := localize(goLocSip)
		err = os.MkdirAll(target, 0755)
		if err != nil {
			return err
		}

		os.RemoveAll(target)
		err = os.Symlink(root, target)
		if err != nil {
			return err
		}
	}

	envStr := 	[]string{
		"GOOS=" + tos,
		"GOARCH=" + tarch,
	}

	err = runGo([]string{"get", pckStr}, envStr, &stdout, &stderr)
	if err != nil {
		fmt.Print(stderr.String())
		return err
	}
	fmt.Print(stdout.String())

	err = runGo(
		[]string{
			"build",
			"-o",
			filepath.Join(localize(locTmp), tos+"-"+tarch),
			pckStr,
		},
		envStr,
		&stdout,
		&stderr,
	)
	if err != nil {
		fmt.Print(stderr.String())
		return err
	}
	fmt.Print(stdout.String())

	return nil
}
