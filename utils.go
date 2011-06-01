package main


import (
	"os"
	"fmt"
	"io/ioutil"
	"path/filepath"
)


func readFile(f string) []byte {
	fc, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Println("error reading file '" + f + "' : " + err.String())
		os.Exit(1)
	}
	return fc
}


func filesInDir(f string, recursive bool) (result []string) {
	result = make([]string, 0, 5)
	return recurse(f, result)
}

func recurse(f string, store []string) ([]string) {
	fi, error := filepath.Glob(f + "/*")
	if error != nil {
		panic("Error while searching for files in dir " + f)
	}
	for _, value := range fi {
		info, _ := os.Stat(value)
		if info.IsDirectory() {
			debug("Recursing into ", value)
			store = recurse(value, store)
		} else {
			store = append(store, value)
		}
	}
	return store
}


func debug(msg string, a ...interface{}) {
	fmt.Printf("DEBUG : "+msg+" : %v\n", a)
}

func printSame(f1, f2 string) {
	fmt.Printf("Same:\r\t%s\r\t%s\r\r", f1, f2)
}

func usage() {
	fmt.Println("Usage: fcmp [OPTION]... FILE [FILE]\nTry `fcmp -h' for more information.\n")
}

func help() {
	fmt.Println("\nfcmp allows to search for duplicate files in folders based on file size and contents.\n")
	usage()
	fmt.Println("FILE can be a regular file or a folder.\n")
	fmt.Println("\tIf two files are given, fcmp compares the two files.")
	fmt.Println("\tIf one file and one folder are given, fcmp compares the file against all the files in the folder.")
	fmt.Println("\tIf one folder is given, fcmp compares all the files in the folder against each other.")
	fmt.Println("\tIf two folders are given, fcmp compares all the files in folder 1 against all the files in folder 2.")
	fmt.Println("\tIf one file is given fcmp prints the usage message and quits\n")

	fmt.Println("Options:")
	fmt.Println("-r\t\tcompare files in subfolder too (ignored if two files are given).")
	fmt.Println("-h\t\tprint this help and exit.")
	fmt.Println("-s\t\tonly compare on file size, do not compare contents (much faster, but of course less reliable).\n")

	fmt.Println("Examples:")
	fmt.Println("Compare two files:\n")
	fmt.Println("\tfcmp foo.jpg bar.jpg\n")
	fmt.Println("Search folder /home/foo for duplicates of file bar.\n")
	fmt.Println("\tfcmp /home/foo bar\n")
	fmt.Println("or\n")
	fmt.Println("\tfcmp bar /home/foo\n")
	fmt.Println("Search folder /home/foo and all its subfolders for duplicates of bar, but hurry up!\n")
	fmt.Println("\tfcmp -r -s /home/foo bar\n")
	fmt.Println("Search folder /home/foo and all its subfolders for duplicates.\n")
	fmt.Println("fcmp -r /home/foo\n")
}
