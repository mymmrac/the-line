package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/tools/godoc/util"
)

func textFilesFromFiles(files []string) (textFiles []string, err error) {
	var isText bool

	for _, file := range files {
		isText, err = isTextFile(file)
		if err != nil {
			return nil, fmt.Errorf("is text: %w", err)
		}

		if isText {
			textFiles = append(textFiles, file)
		}
	}

	return textFiles, nil
}

func isTextFile(filename string) (bool, error) {
	//nolint:gosec
	file, err := os.Open(filename)
	if err != nil {
		return false, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var buf [1024]byte
	n, err := file.Read(buf[:])
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false, nil
		}
		return false, fmt.Errorf("read file: %w", err)
	}

	return util.IsText(buf[:n]), nil
}
