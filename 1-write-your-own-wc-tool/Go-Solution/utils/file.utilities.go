package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func GetFileInformation(filename string) (string, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}

	// prepate the return message
	message := fmt.Sprint("\nInformation",
		"\n File Name: ", fileInfo.Name(),
		"\n Permissions: ", fileInfo.Mode(),
		"\n Last Modified: ", fileInfo.ModTime(),
	)

	// return
	return message, err
}

func GetFileSize(filename string) (string, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	// prepate the return message
	message := fmt.Sprintf("\n Size (in bytes): %d", fileInfo.Size())

	// return
	return message, err
}

func GetFileLines(filename string) (string, error) {
	// initiate file-hanndle to read from
	file, err := os.Open(filename)

	// check if file-handle was initiated correctly
	if err != nil {
		log.Fatal(err)
	}

	// initiate scanner from file-handle
	fileScanner := bufio.NewScanner(file)

	// tell the scanner to split by lines
	fileScanner.Split(bufio.ScanLines)

	// initiate the counter & looping through lines
	lineCount := 0
	for fileScanner.Scan() {
		lineCount += 1
	}

	// make sure to close the file-hanlde upon return
	defer file.Close()

	// prepate the return message
	message := fmt.Sprintf("\n Total lines: %d", lineCount)

	// return
	return message, err
}

func GetFileCharacters(filename string) (string, error) {
	// initiate the file-hanndle to read from
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	// initiate the scanner from file-handle
	fileScanner := bufio.NewScanner(file)

	// tell the scanner to split the chars
	fileScanner.Split(bufio.ScanRunes)

	// initiate the counter & looping through the runes(ASCII) / character
	charCount := 0
	for fileScanner.Scan() {
		charCount += 1
	}

	// make sure to close the file-handle upon return
	defer file.Close()

	// prepate the return message
	message := fmt.Sprintf("\n Total characters: %d", charCount)

	// return
	return message, err
}

func GetFileWords(filename string) (string, error) {
	// initiate file-handle to read from
	file, err := os.Open(filename)

	// check if file-handle was initiated correctly
	if err != nil {
		log.Fatal(err)
	}

	// initiate scanner from file-handle
	fileScanner := bufio.NewScanner(file)

	// tell the scanner to split by words
	fileScanner.Split(bufio.ScanWords)

	// initiate counter & looping through words
	wordCount := 0
	for fileScanner.Scan() {
		wordCount += 1
	}

	// check if there was error while reading words from file
	if err := fileScanner.Err(); err != nil {
		log.Fatal(err)
	}

	// make to close file-handle upon return
	defer file.Close()

	// prepate the return message
	message := fmt.Sprintf("\n Total words: %d", wordCount)

	// return
	return message, err
}
