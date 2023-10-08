# WC Tool

The Unix command line tools are a great metaphor for good software engineering and they follow the Unix Philosophies of:

Writing simple parts connected by clean interfaces - each tool does just one thing and provides a simple CLI that handles text input from either files or file streams.
Design programs to be connected to other programs - each tool can be easily connected to other tools to create incredibly powerful compositions.
Following these philosophies has made the simple unix command line tools some of the most widely used software engineering tools - allowing us to create very complex text data processing pipelines from simple command line tools. There’s even a Coursera course on [Linux and Bash for Data Engineering](https://gb.coursera.org/learn/linux-and-bash-for-data-engineering-duke).

You can read more about the Unix Philosophy in the excellent book [The Art of Unix Programming](http://www.catb.org/~esr/writings/taoup/html/).

## Installation

The functional requirements for wc are concisely described by it’s man page - give it a go in your local terminal now:

```bash
man wc
```

## Usage

```go
$ go run .

WC Arguments:
        -f FILENAME, --filename FILENAME 
                                                The input file, or standard input (if no file is specified) to the standard output.
        -c, --bytes         The number of bytes in each input file is written to the standard output.
        -l, --lines         The number of lines in each input file is written to the standard output.
        -m, --characters    The number of characters in each input file is written to the standard output.
        -w, --words         The number of words in each input file is written to the standard output.

 Enter the filename here: 
test.txt

 Want to get size in bytes: [y/n]


 Want to get total number of lines in the file: [y/n]


 Want to get total number of characters in the file: [y/n]


 Want to get total number of words in the file: [y/n]


Information 
 File Name:  test.txt
 Permissions:  -rw-r--r--
 Last Modified:  2023-10-07 20:58:50.640676903 +0530 IST

 Size (in bytes):  335041

 Total lines:  7144

 Total characters:  332145

 Total words:  58164
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[GNU General Public License v3.0](https://github.com/abhishekpatel946/coding-challeges/blob/ba3743f975ac8f2b8ad5757813a729336d131955/LICENSE)