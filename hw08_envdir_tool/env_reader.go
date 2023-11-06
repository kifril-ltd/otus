package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path"
	"strings"
	"sync"
)

var ErrNotADirectory = errors.New("provided path is not a directory")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns a map of environment variables.
// Variables represented as files where the filename is the name of the variable,
// and the file's first line is the value.
func ReadDir(dir string) (Environment, error) {
	if !isDir(dir) {
		return nil, ErrNotADirectory
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	wg.Add(len(files))
	for _, file := range files {
		go func(file os.DirEntry) {
			defer wg.Done()

			if !isEnvFile(file) {
				return
			}

			key, value, err := readEnvFile(dir, file)
			if err != nil {
				return
			}

			mu.Lock()
			env[key] = value
			mu.Unlock()
		}(file)
	}

	wg.Wait()

	return env, nil
}

func isDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func isEnvFile(file os.DirEntry) bool {
	fileName := file.Name()
	return !file.IsDir() && !strings.Contains(fileName, "=")
}

func readEnvFile(dir string, file os.DirEntry) (string, EnvValue, error) {
	filePath := path.Join(dir, file.Name())
	f, err := os.Open(filePath)
	if err != nil {
		return "", EnvValue{}, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return "", EnvValue{}, err
	}

	if info.Size() == 0 {
		return file.Name(), EnvValue{Value: "", NeedRemove: true}, nil
	}

	br := bufio.NewReader(f)
	line, _, err := br.ReadLine()
	if err != nil {
		return "", EnvValue{}, err
	}

	value := normalizeValue(line)

	return file.Name(), EnvValue{Value: value, NeedRemove: false}, nil
}

func normalizeValue(value []byte) string {
	normalized := bytes.ReplaceAll(value, []byte("\x00"), []byte("\n"))
	normalized = bytes.TrimRight(normalized, " \t\n")

	return string(normalized)
}
