package filemanager

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type FileManager struct {
	InputFilePath  string
	OutputFilePath string
}

func (fm *FileManager) ReadLines() ([]string, error) {
	// open a file and return a slice of strings
	file, err := os.Open(fm.InputFilePath)
	if err != nil {
		fmt.Println("Could not open file!")
		return nil, err
	}

	// when opening a file succeeded
	// defer Close operation once the ReadLines() func finishes execution
	defer file.Close()

	// utility for reading IO
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	var lines []string
	for scanner.Scan() { // True until there are no more lines on the file
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

// OpenWriter opens a path for writing as a io.WriteCloser interface
func (fm *FileManager) OpenWriter() (io.WriteCloser, error) {
	f, err := os.Create(fm.OutputFilePath)
	if err != nil {
		return nil, fmt.Errorf("filemanager: opening %q: %w", fm.OutputFilePath, err)
	}
	return f, nil
}

func (fm *FileManager) WriteResult(data any) error {
	// create a file with the os.Create() method
	file, err := os.Create(fm.OutputFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func New(inputPath, outputPath string) FileManager {
	return FileManager{
		InputFilePath:  inputPath,
		OutputFilePath: outputPath,
	}
}
