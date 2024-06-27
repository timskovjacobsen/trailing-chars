# Highlight trailing characters (e.g. whitespaces)

Simple CLI to highlight trailing characters (whitespaces by default) in source files.

```shell
❯ trailing --help
trailing - highlight trailing characters in files

Usage: trailing [OPTIONS] <FILENAMES...>

OPTIONS:
  --chars string
    	The trailing chars (in any combination) to detect at line ends
    	(default: all whitespace)

EXAMPLES:
  check trailing whitespace in main.go & main.py:

  	trailing main.go main.py

  check trailing characters 'abc' in main.go:

  	trailing --chars 'abc' main.go
```
