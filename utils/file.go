package utils

import (
	"os"
	"path/filepath"
	"strings"
	//"github.com/aywfelix/felixgo/container"
)

var (
	Separator = string(filepath.Separator)
)

type file struct{}

func (f file) Mkdir(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (f file) Create(path string) (*os.File, error) {
	dir := f.Dir(path)
	if !f.Exists(dir) {
		if err := f.Mkdir(dir); err != nil {
			return nil, err
		}
	}
	return os.Create(dir)
}

func (f file) Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

func (f file) Join(paths ...string) string {
	var s string
	// for _, path := range paths {
	// 	if s != "" {
	// 		s += Separator
	// 	}
	// 	s += strings.TrimRight(path, container.DefaultTrimChars)
	// }
	return s
}

func (f file) IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func (f file) IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func (f file) Pwd() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

func (f file) Dir(path string) string {
	if path == "" {
		return filepath.Dir(f.RealPath(path))
	}
	return filepath.Dir(path)
}

func (f file) RealPath(path string) string {
	p, err := filepath.Abs(path)
	if err != nil || !f.Exists(p) {
		return ""
	}
	return p
}

func (f file) DirNames(path string) ([]string, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	list, err := fp.Readdirnames(-1)
	defer fp.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (f file) Basename(path string) string {
	return filepath.Base(path)
}

func (f file) Name(path string) string {
	base := filepath.Base(path)
	if i := strings.LastIndexByte(base, '.'); i != -1 {
		return base[:i]
	}
	return base
}

func (f file) IsPathEmpty(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return true
	}
	if stat.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			return true
		}
		defer file.Close()
		names, err := file.Readdirnames(-1)
		if err != nil {
			return true
		}
		return len(names) == 0
	} else {
		return stat.Size() == 0
	}
}

func (f file) Ext(path string) string {
	ext := filepath.Ext(path)
	if p := strings.IndexByte(ext, '?'); p != -1 {
		ext = ext[0:p]
	}
	return ext
}

func (f file) ExtName(path string) string {
	return strings.TrimLeft(f.Ext(path), ".")
}

func (f file) GetCurrentDir() string {
	dir, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic("get current directory error")
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//========================================================================
var File file
