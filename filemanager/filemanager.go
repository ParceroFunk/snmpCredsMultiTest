package filemanager

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type FileManager struct {
	InputFilePath  string
	OutputFilePath string
}

func (fm FileManager) ReadLines() ([]string, error) {
	// open a file and return a slice of strings
	file, err := os.Open(fm.InputFilePath)
	if err != nil {
		fmt.Println("Could not open file!")
		fmt.Println(err)
		return nil, errors.New("failed to open specified file")
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
		fmt.Println(err)
		return nil, errors.New("error while scanning the file")
	}

	file.Close()
	return lines, nil
}

func (fm FileManager) WriteResult(data any) error {
	// create a file with the os.Create() method
	file, err := os.Create(fm.OutputFilePath)
	if err != nil {
		fmt.Println(err)
		return errors.New("error while creating the file")
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		fmt.Println(err)
		return errors.New("error while creating the file")
	}

	return nil
}

func New(inputPath, outputPath string) FileManager {
	return FileManager{
		InputFilePath:  inputPath,
		OutputFilePath: outputPath,
	}
}
