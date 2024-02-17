package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %q not exists", dir)
	}

	env := Environment(make(map[string]EnvValue))

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if d.IsDir() || strings.ContainsRune(d.Name(), '=') {
			return nil
		}

		file := filepath.Join(dir, d.Name())

		val, err := readFile(file)
		if err != nil {
			return err
		}

		env[d.Name()] = EnvValue{
			Value:      val,
			NeedRemove: val == "",
		}

		return nil
	})

	return env, err
}

func readFile(file string) (string, error) {
	f, fErrOpen := os.Open(file)
	if fErrOpen != nil {
		return "", fErrOpen
	}

	defer f.Close()

	r := bufio.NewReader(f)

	val, fErrRead := r.ReadBytes('\n')
	if fErrRead != nil {
		if errors.Is(fErrRead, io.EOF) {
			return convertValue(val), nil
		}

		return "", fErrRead
	}

	val = val[0 : len(val)-1]

	return convertValue(val), nil
}

func convertValue(val []byte) string {
	if len(val) > 0 {
		val = bytes.ReplaceAll(val, []byte{0x00}, []byte{'\n'})
		val = bytes.TrimRight(val, " \t")
	}

	return string(val)
}
