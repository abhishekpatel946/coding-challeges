package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/abhishekpatel946/1-write-your-own-wc-tool/utils"
)

// Configuration
const Message string = `
WC Arguments:
	-f FILENAME, --filename FILENAME 
						The input file, or standard input (if no file is specified) to the standard output.
	-c, --bytes         The number of bytes in each input file is written to the standard output.
	-l, --lines         The number of lines in each input file is written to the standard output.
	-m, --characters    The number of characters in each input file is written to the standard output.
	-w, --words         The number of words in each input file is written to the standard output.
`

// Throwing the custom error while parsing arguments
var ErrParse = errors.New("invalid arguments passed to command line, please specify the [y/n] argument accepted only")

// Custom type cast fn for string to boolean conversion based on [y/n] argument only
func ConvertStrToBool(s string) (bool, error) {
	switch s {
	case "y":
		return true, nil
	case "":
		return true, nil
	case "n":
		return false, nil
	}
	return false, ErrParse
}

func wc(filename string, _bytes bool, _lines bool, _characters bool, _words bool) {
	res, err := utils.GetFileInformation(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)

	// get the size of the file
	if _bytes {
		res, err := utils.GetFileSize(filename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	}

	// get the total no of lines in the file
	if _lines {
		res, err := utils.GetFileLines(filename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	}

	// get the total no of chars in the file
	if _characters {
		res, err := utils.GetFileCharacters(filename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	}

	// get the total no of words in the file
	if _words {
		res, err := utils.GetFileWords(filename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	}
}

func checkFileExists(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	} else {
		return false, err
	}
}

// main function
func main() {

	// Print arguments required from command line
	fmt.Print(Message)

	// Take the list of arguments from command line & type-case the string argument to boolean
	var filename string
	var _bytes string
	var _lines string
	var _characters string
	var _words string

	fmt.Println("\n Enter the filename here: ")
	fmt.Scanln(&filename)

	// check the filename is exist or not
	isFileExist, err := checkFileExists(filename)
	if err != nil {
		log.Fatal(err)
	}

	if isFileExist {
		fmt.Println("\n Want to get size in bytes: [y/n]")
		fmt.Scanln(&_bytes)
		_bytesBool, err := ConvertStrToBool(_bytes)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\n Want to get total number of lines in the file: [y/n]")
		fmt.Scanln(&_lines)
		_linesBool, err := ConvertStrToBool(_lines)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\n Want to get total number of characters in the file: [y/n]")
		fmt.Scanln(&_characters)
		_charactersBool, err := ConvertStrToBool(_characters)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\n Want to get total number of words in the file: [y/n]")
		fmt.Scanln(&_words)
		_wordsBool, err := ConvertStrToBool(_words)
		if err != nil {
			log.Fatal(err)
		}

		// Compute the file operations
		wc(filename, _bytesBool, _linesBool, _charactersBool, _wordsBool)
	}
}
