package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
)
 

func getFileInformation(filename string) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Log the basic information about the file
	fmt.Print("\nInformation \n")
	fmt.Println(" File Name: ", fileInfo.Name())
	fmt.Println(" Permissions: ", fileInfo.Mode())
	fmt.Println(" Last Modified: ", fileInfo.ModTime())
}

func getFileSize(filename string) (string, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	// return
	return "\n Size (in bytes): " + string(rune(fileInfo.Size())), nil
}

func getFileLines(filename string) (string, error) {
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

	// return
	return "\n Total lines: " + string(lineCount), nil
}

func getFileCharacters(filename string) (string, error) {
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

	return "\n Total characters: " + string(rune(charCount)), nil
}

func getFileWords(filename string) (string, error) {
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

	return "\n Total words: " + string(rune(wordCount)), nil
}
