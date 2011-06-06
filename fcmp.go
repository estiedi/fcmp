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

func (cmp *ComparableFile) CompareSize(path string, f *os.FileInfo) bool {
	result := cmp.info.Size == f.Size
	if *verbose {
		fmt.Printf("Size of %s == size of %s is %t\n", cmp.name, path, result)
	}
	return result
}

func (cmp *ComparableFile) CompareContents(path string) bool {
	fc1 := readFile(cmp.name)
	fc2 := readFile(path)
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

func (cmp *ComparableFile) Compare(path string, f *os.FileInfo) {
	same := cmp.CompareSize(path, f)
	if *sizeOnly && same {
		printSame(cmp.name, path)
	}else{
		if same {
			same = cmp.CompareContents(path)
			if same {
				printSame(cmp.name, path)			
			}
		}
	}
}

func (cmp *ComparableFile) VisitDir(path string, f *os.FileInfo) bool{
	return true;
}

func (cmp *ComparableFile) VisitFile(path string, f *os.FileInfo){
	if *verbose {	
		fmt.Printf("Compare %s  with %s\n", cmp.name, path)
	}
	cmp.Compare(path, f)
}

func createFileSet(filename string, info *os.FileInfo) []*ComparableFile {
	if info.IsDirectory() {
		filesInDir := filesInDir(filename, true)
		fileset := make([]*ComparableFile, len(filesInDir))
		for idx, value := range filesInDir {
			fileset[idx] = NewComparableFile(value)
		}
		return fileset
	}
	return []*ComparableFile{NewComparableFile(filename)}
}

var verbose = flag.Bool("v", false, "print lots of messages about what'sgoing on")
var needHelp = flag.Bool("h", false, "print help and exit")
var sizeOnly = flag.Bool("s", false, "only compare based on file size")

func main() {
	flag.Parse()
	if *needHelp {
		help()
		os.Exit(0)
	}
	fileSet, target := initialize()
	for _, value := range fileSet {
		if *verbose {
			fmt.Printf("Processing %s\n", value.name)
		}
		filepath.Walk(target, value, nil)	
	}
}

func initialize() ([]*ComparableFile, string) {
	args := flag.NArg()
	fn1 := flag.Arg(0)
	fn2 := flag.Arg(1)
	if *verbose {
		fmt.Printf("Args: %s : %s\n", fn1, fn2)
	}
	// if no args, use the current directory
	if args == 0 {
		fn1,_ = os.Getwd()
		fn2 = fn1
	}
	fi1, _ := os.Stat(fn1)
	isDir := fi1.IsDirectory()
	// If  fn1 is a file compare it against all the files in the same dir, if it's a dir compare all files in dir against each other
	if args == 1 {
		if isDir {
			fn2 = fn1
		}else{
			fn2, _ = filepath.Split(fn1)		
		}
	}
	filesToCompare := createFileSet(fn1, fi1) 
	if *verbose {
		fmt.Printf("Got %d files to compare with %s\n", len(filesToCompare), fn2)
	}	
	return filesToCompare, fn2
}


