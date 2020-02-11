package main

import (
	"fmt"
	"github.com/tjenkins13/tarfs/tinfo"
	"github.com/tjenkins13/tarfs/tree"
	"os"
	"sort"
)

type TarFile = tinfo.TarFile
type Node = tree.Node

type OpenStruct struct {
	tname   string   //absolute path to tar file
	mloc    string   //absolute path to where you are mounting tar file
	fname   string   //name of file with tar file you want
	blockno int64    //block number you are at in file
	offset  int64    //offset within block that you are at
	fd      *os.File //file descriptor for tar file
}

//For Sorting List of Tarfile headers by length of name
//--guarantees that root header is first
type ByNameLen []TarFile

func (a ByNameLen) Len() int           { return len(a) }
func (a ByNameLen) Less(i, j int) bool { return len(a[i].Name) < len(a[j].Name) }
func (a ByNameLen) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//checks if we have read in a byte array of 512 0's, which is what the tar file ends in
func allZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

//Opens tar file-- will have other function to open file withing tar file
func OpenTar(filename string) int {
	var tfile []TarFile
	var t TarFile
	var block int64 = 0
	var head Node
	sread := 0
	dat, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	b1 := make([]byte, 512)
	_, err2 := dat.Read(b1)
	if err2 != nil {
		panic(err2)
	}
	//fmt.Println(string(b1[:100]))
	//fmt.Println(b1)
	for !allZero(b1) {
		if sread == 0 {
			t.Create(b1, block)
			//fmt.Println(t.Name)
			tfile = append(tfile, t)
			block++
			if tfile[len(tfile)-1].Typeflag != '5' { //if it's not a directory we're about to read some data
				sread = int((tfile[len(tfile)-1].Size-1)/512) + 1
			}
		} else {
			sread--
		}
		dat.Read(b1)
		if allZero(b1) { //tar file ends with two empty blocks and we found one
			dat.Read(b1) //get second block- loop checks if its empty
			block++
		}
	}
	sort.Sort(ByNameLen(tfile)) // sort by length of file name for insertion into tree
	//get tree of filesystem and put it in array
	head = tree.Create(tfile[0])
	//fmt.Println(head.Data.Name)
	for _, file := range tfile[1:] {
		head.Insert(file)
		//fmt.Println(file.Name)
	}
	head.Print()
	return len(tfile) //for now, will be change to return element of array

}

func CloseTar(fd int) {

}

func main() {
	//fmt.Println("Hello World")
	fmt.Println(OpenTar("data/x.tar"))
}
