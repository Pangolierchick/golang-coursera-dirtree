package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func printSpaces(depth int) string {
	return strings.Repeat(" ", depth * 4)
}

func remove(slice []os.FileInfo, ind int) []os.FileInfo {
	return append(slice[:ind], slice[ind + 1:]...)
}

func printDir(outf *os.File, path string, printFiles bool, depth int, pref string) error {
	files, err := ioutil.ReadDir(path)

	if printFiles == false {
		for i := 0; i < len(files) ; {
			if files[i].IsDir() == false {
				files = remove(files, i)
			} else {
				i++;
			}
		}
	}

	if err != nil {
		return fmt.Errorf("Couldn't read dirs")
	}

	// sort.Stable(sortFiles(files))

	for i, file := range files {
		var symb string
		var newpref string
		// println("[DBG] files ", len(files))
		if i == len(files) - 1 {
			symb = "└───"
			newpref = "    "
		} else {
			symb = "├───"
			newpref = "│   "
		}

		// fmt.Printf("Curr file %d: %s\n", i, file.Name())

		
		// println("[DBG] Curr file: ", file.Name(), file.IsDir())
		if printFiles && !file.IsDir() {
			fmt.Fprintf(outf, "%s%s%s (%db)\n", pref, symb, file.Name(), file.Size())
		} else if file.IsDir() {
			fmt.Fprintf(outf, "%s%s%s\n", pref, symb, file.Name())
			printDir(outf, path + string(os.PathSeparator) + file.Name(), printFiles, depth + 1, pref + newpref)
		}
	}

	return nil
}

func dirTree(outf *os.File, path string, printFiles bool) error {
	return printDir(outf, path, printFiles, 0, "")
}

func main() {
	out := os.Stdout

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)

	if err != nil {
		panic(err.Error())
	}
}
