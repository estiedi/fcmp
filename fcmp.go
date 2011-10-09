/*
This file is part of fcmp.

fcmp is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
any later version.

fcmp is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with fcmp.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"os"
	"path/filepath"
	"fmt"
	"regexp"
	"crypto/md5"
	"io"
	"flag"
	"time"
)

var pattern = flag.String("p", ".*", "pattern to select files to compare (e.g. \"^.*.jpg$|^.*JPG$\" selects all jpg files)")

type Md5File struct {
	Path string
	Md5  string
}

type FileVisitor struct {
	out    chan *Md5File
	regexp *regexp.Regexp
}


func main() {
	flag.Parse()
	filelist := make([]Md5File, 0, 100)
	out := make(chan *Md5File)
	errors := make(chan os.Error)
	root := "./"
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}
	visitor := NewFileVisitor(out, *pattern)
	go func() {
		count := 0
		for {
			select {
			case err := <-errors:
				fmt.Println("no error expected, found: %s", err)
			case result := <-visitor.out:
				filelist = append(filelist, *result)
				count++
				fmt.Print("*")
				if count % 80 == 0 {
					fmt.Printf(": %d\n", count)
				}
			}
		}
	}()
	fmt.Println("Calculating md5 sum for files:")
	startTime := time.Seconds()
	filepath.Walk(root, visitor, errors)
	stopTime := time.Seconds()
	fmt.Printf("\nFinished calculating md5 sums. Files : %d in %d sec.\n", len(filelist), stopTime - startTime)
	fmt.Println("Start sorting...")
	bubblesort(filelist)
	for idx := 0; idx < len(filelist)-1; idx++ {
		if filelist[idx].Md5 == filelist[idx+1].Md5 {
			fmt.Printf("%s : %s\n", filelist[idx].Path, filelist[idx+1].Path)
		}
	}
	stopTime = time.Seconds()
	fmt.Printf("Total time : %d sec.\n", stopTime - startTime)
}


func (FileVisitor) VisitDir(path string, f *os.FileInfo) bool {
	return true
}
func (visitor FileVisitor) VisitFile(path string, f *os.FileInfo) {
	if visitor.regexp.MatchString(path) {
		md5 := calcMd5(path)
		visitor.out <- &Md5File{path, md5}
	}
}

func calcMd5(path string) string {
	md5h := md5.New()
	md5h.Reset()
	inFile, _ := os.Open(path)
	defer inFile.Close()
	io.Copy(md5h, inFile)
	sum := fmt.Sprintf("%x", md5h.Sum())
	return sum
}

func bubblesort(a []Md5File) {
	for itemCount := len(a) - 1; ; itemCount-- {
		hasChanged := false
		for index := 0; index < itemCount; index++ {
			if a[index].Md5 > a[index+1].Md5 {
				a[index], a[index+1] = a[index+1], a[index]
				hasChanged = true
			}
		}
		if hasChanged == false {
			break
		}
	}
}

func NewFileVisitor(out chan *Md5File, pattern string) *FileVisitor {
	regexp := regexp.MustCompile(pattern)
	return &FileVisitor{out, regexp}
}
