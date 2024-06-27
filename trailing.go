package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"unicode"
)

func main() {
	flagSet := flag.NewFlagSet("trailing", flag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "trailing - highlight trailing characters in files\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <FILENAMES...>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "  --chars string\n")
		fmt.Fprintf(os.Stderr, "    \tThe trailing chars (in any combination) to "+
			"detect at line ends\n")
		fmt.Fprintf(os.Stderr, "    \t(default: all whitespace)\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "EXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "  check trailing whitespace in main.go & main.py:\n\n")
		fmt.Fprintf(os.Stderr, "  \ttrailing main.go main.py\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  check trailing characters 'abc' in main.go:\n\n")
		fmt.Fprintf(os.Stderr, "  \ttrailing --chars 'abc' main.go\n\n")
	}
	chars := flagSet.String("chars", "", "The trailing chars to detect at line ends (default: all whitespace)")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		flagSet.Usage()
		return
	}

	if flagSet.NArg() < 1 || flagSet.Args()[0] == "help" {
		flagSet.Usage()
		return
	}
	filenames := flagSet.Args()

	for _, filename := range filenames {
		checkFile(filename, *chars)
	}
}

func checkFile(filename string, chars string) bool {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open file: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineno := 0
	for scanner.Scan() {
		lineno++
		line := scanner.Bytes()
		hlLine := highlightTrailing(line, chars)
		if hlLine != nil {
			fmt.Printf("%s:%d\n", filename, lineno)
			fmt.Println(string(hlLine))

			fmt.Println()
		}
	}
	return false
}

func lineParts(line []byte) ([]byte, []byte) {
	var lineStart []byte
	var eol []byte
	if bytes.HasSuffix(line, []byte("\n")) {
		lineStart = line[:len(line)-1]
		eol = []byte("\n")
	} else if bytes.HasSuffix(line, []byte("\r\n")) {
		lineStart = line[:len(line)-2]
		eol = []byte("\r\n")
	} else {
		lineStart = line
		eol = []byte("")
	}
	return lineStart, eol
}

func highlightTrailing(line []byte, chars string) []byte {
	lineStart, eol := lineParts(line)

	var stripped []byte
	if len(chars) == 0 {
		// Default chars -> all whitepace
		stripped = bytes.TrimRightFunc(lineStart, unicode.IsSpace)
	} else {
		stripped = bytes.TrimRight(lineStart, chars)
	}
	var newline []byte
	if bytes.Equal(line, stripped) {
		return nil
	} else {
		RED := []byte("\033[41m")
		RESET := []byte("\033[0m")
		trailingChars := lineStart[len(stripped):]
		newline = append(newline, stripped...)
		newline = append(newline, RED...)
		newline = append(newline, trailingChars...)
		newline = append(newline, RESET...)
		newline = append(newline, eol...)
	}
	return newline
}
