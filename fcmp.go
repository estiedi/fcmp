package main


import (
	"os"
	"flag"
	"path/filepath"
	"fmt"
)

type ComparableFile struct{
	name string
	info *os.FileInfo
}

func NewComparableFile(filename string) (*ComparableFile) {
	fi,err := os.Stat(filename)
	if err != nil {
		fmt.Println(err.String())
		os.Exit(1)	
	}
	return &ComparableFile{name: filename, info: fi}
}

func (cmp *ComparableFile) CompareSize(tgt *ComparableFile) bool {
	return cmp.info.Size == tgt.info.Size
}

func (cmp *ComparableFile) CompareContents(tgt *ComparableFile) bool {
	fc1 := readFile(cmp.name)
	fc2 := readFile(tgt.name)
	la := len(fc1)
	lb := len(fc2)
	smallest := lb
	if la < lb {
		smallest = la
	}
	for i := 0; i < smallest; i++ {
		if fc1[i] != fc2[i] {
			return false
		}
	}
	return true
}

func (cmp *ComparableFile) Compare( tgt *ComparableFile) {
	same := cmp.CompareSize(tgt)
	if *sizeOnly && same {
		printSame(cmp.name, tgt.name)
	}else{
		if same {
			same = cmp.CompareContents(tgt)
			if same {
				printSame(cmp.name, tgt.name)			
			}
		}
	}
}

/*
Parses and validates the program invocation arguments.
If arguments are valid, it constructs two sets of comparable files. 
Files in set 1 will be compared to files in set 2.
If an argument is invalid, prints an error message and exits.
*/
func initialize() ([]*ComparableFile, []*ComparableFile) {
	args := flag.NArg()
	fn1 := flag.Arg(0)
	fn2 := flag.Arg(1)
	fi1, _ := os.Stat(fn1)

	switch args {
// if no args, use the current directory
	case 0:
		usage()
		os.Exit(0)
	case 1:
// Don't care about this. If it-s a file compare it against itself, if it's a dir compare all files in dir against each other
		isDir := fi1.IsDirectory()
		if isDir {
			fn2 = fn1
		} else {
			fn2, _ = filepath.Split(fn1)
		}
	case 2:
		// do nothing, just use the first two arguments
	}
/*
	filesInDir1 := filesInDir(fn1, true)

	filesInDir2 := filesInDir(fn2, true)

	fileset1 := make([]*ComparableFile, len(filesInDir1))
	fileset2 := make([]*ComparableFile, len(filesInDir2))
	for idx, value := range filesInDir1 {
		fileset1[idx] = NewComparableFile(value)
	}
	for idx, value := range filesInDir2 {
		fmt.Printf("Index %d : value: %s\n", idx, value)
		fileset2[idx] = NewComparableFile(value) //index out of range error
	}
*/	
	fileset1 := createFileSet(fn1)
	fileset2 := createFileSet(fn2)
	return fileset1, fileset2
}

func createFileSet(filename string) []*ComparableFile {
	filesInDir := filesInDir(filename, true)
	fileset := make([]*ComparableFile, len(filesInDir))
	for idx, value := range filesInDir {
		fileset[idx] = NewComparableFile(value)
	}
	return fileset
}

func compareFileToDir(src *ComparableFile, tgt []*ComparableFile) {
	for _, value := range tgt {
		src.Compare(value)
	}
}

func compareDirToDir(src, tgt []*ComparableFile) {
	for _, srcFile := range src {
		compareFileToDir(srcFile, tgt)	
	}
}



var recursive = flag.Bool("r", false, "search subfolders too")
var needHelp = flag.Bool("h", false, "print help and exit")
var sizeOnly = flag.Bool("s", false, "only compare based on file size")

func main() {
	flag.Parse()
	if *needHelp {
		help()
		os.Exit(0)
	}
	srcFileSet, tgtFileSet := initialize()
	switch {
	case len(srcFileSet) == 1:
		compareFileToDir(srcFileSet[0], tgtFileSet)
	case len(srcFileSet) > 1:
		compareDirToDir(srcFileSet, tgtFileSet)

	}
}


