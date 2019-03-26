package index

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path"
	"strconv"
)

var chunkSize int64 = 4096

type FileSystem struct {
	index *Index
	Files map[string]File
}

type File struct {
	Path   string
	Size   int64
	Chunks []int
}

func (f *FileSystem) List() map[string]File {
	if f.Files == nil {
		f.Files = map[string]File{}
	}

	// TODO: sync with other nodes

	return f.Files
}

func (f *FileSystem) internalID(p string) string {
	b := sha256.Sum256([]byte(path.Clean(path.Join("/", p))))
	return hex.EncodeToString(b[:])
}

func (f *FileSystem) internalPath(p string) string {
	return path.Join("./.tmp/fs", f.internalID(p)) + "/"
}

func (f *FileSystem) publicPath(p string) string {
	return path.Clean(path.Join("/", p))
}

func (f *FileSystem) Add(p string, c io.Reader) error {
	p = f.publicPath(p)
	intPath := f.internalPath(p)

	s := chunkSize
	os.MkdirAll(intPath, os.ModePerm)
	chunkCount := 0

	newFile := File{
		Path:   p,
		Size:   0,
		Chunks: []int{},
	}

	for s == chunkSize {
		fil, err := os.Create(intPath + "/" + strconv.Itoa(chunkCount) + ".chunk")
		newFile.Chunks = append(newFile.Chunks, chunkCount)
		chunkCount++

		if err != nil {
			return err
		}
		defer fil.Close()
		s, err = io.CopyN(fil, c, chunkSize)
		newFile.Size = newFile.Size + s

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	if f.Files == nil {
		f.Files = map[string]File{}
	}
	f.Files[newFile.Path] = newFile
	return nil
}

func (f *FileSystem) Get(p string, c io.Writer) error {
	p = f.publicPath(p)
	intPath := f.internalPath(p)

	cf, ex := f.Files[p]
	if !ex {
		return errors.New("File does not exist.")
	}

	// All is local (for first version read from local disk)
	totalChunks := int(cf.Size / chunkSize)
	chunkLen := len(cf.Chunks) - 1
	if chunkLen == totalChunks {
		chunkCount := 0
		for chunkCount < totalChunks {
			fil, err := os.Open(intPath + "/" + strconv.Itoa(chunkCount) + ".chunk")
			chunkCount++

			if err != nil {
				return err
			}
			defer fil.Close()
			_, err = io.CopyN(c, fil, chunkSize)

			if err == io.EOF {
				break
			}

			if err != nil {
				return err
			}
		}
		return nil
	}

	// TODO: allow for fetching from multiple sources
	return errors.New("File does not exist locally (not yet implemented).")
}

func NewFileSystem(i *Index) FileSystem {
	fs := FileSystem{
		index: i,
	}

	return fs
}
