package main

import (
	"testing"
	"os"
//	"fmt"
)


func TestCompareSizeSame(t *testing.T) {
	file1 := "./testdata/pierre.jpg"
	file2 := "./testdata/pierre2.jpg"
	fi1, _ := os.Stat(file1)
	fi2, _ := os.Stat(file2)

	result := compareSize(fi1, fi2)
	if !result {
		t.Errorf("Expected files to be the same!")
	}
}

func TestCompareSizeDifferent(t *testing.T) {
	file1 := "./testdata/pierre.jpg"
	file2 := "./testdata/changedExif.jpg"
	fi1, _ := os.Stat(file1)
	fi2, _ := os.Stat(file2)

	result := compareSize(fi1, fi2)
	if result {
		t.Errorf("Expected files to be different!")
	}
}

func TestCompareContentsSame(t *testing.T) {
	file1 := "./testdata/PatternsInJava_2002.chm"
	file2 := "./testdata/Wiley-Patterns.In.Java-A.Catalog.Of.Reusable.Design.Patterns.Illustrated.With.Uml-Vol.1.2E.chm"
	result := compareContents(file1, file2)
	if !result {
		t.Errorf("Expected files to be the same!")
	}

}
/*
func TestFileFunctions(t *testing.T) {
	file1 := "./testdatafoobar/PatternsInJava_2002.chm"
	fi1, err := os.Stat(file1)
	if err != nil {
		fmt.Println("stat error:"+err.String())
		return
	}
	if fi1.IsDirectory() {
		fmt.Println("Is a directory")
	} else {
		fmt.Println("Is NOT a directory")
		if fi1.IsRegular() {
			fmt.Println("It's a regular file")
		} else {
			fmt.Println("and not even a regular file")
		}
	}
}
*/

