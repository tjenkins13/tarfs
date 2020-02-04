package main

import "fmt"

//import "bufio"
//import "io"
//import "io/ioutil"
import "os"
import "strconv"
import "strings"

//import "time"

/*
type TarFile struct{
name [100]byte
mode [8]byte
uid [8]byte
gid [8]byte
size [12]byte
mtime [12]byte
chksum [8]byte
typeflag byte
linkname [100]byte
magic [8]byte
uname [32]byte
gname [32]byte
devmajor [8]byte
devminor [8]byte
prefix [167]byte
}
*/

type TarFile struct {
	name     []byte
	mode     []byte
	uid      []byte
	gid      []byte
	size     []byte
	mtime    []byte
	chksum   []byte
	typeflag byte
	linkname []byte
	magic    []byte
	uname    []byte
	gname    []byte
	devmajor []byte
	devminor []byte
	prefix   []byte
}

type TarFile2 struct {
	name     string
	mode     string
	uid      int64
	gid      int64
	size     int64
	mtime    int64
	chksum   int64
	typeflag byte
	linkname string
	magic    string
	uname    string
	gname    string
	devmajor int64
	devminor int64
	prefix   string
}

func (tFile *TarFile) Create(b1 []byte) {
	//Allocate memory for parts of struct
	tFile.name = make([]byte, 100)
	tFile.mode = make([]byte, 8)
	tFile.uid = make([]byte, 8)
	tFile.gid = make([]byte, 8)
	tFile.size = make([]byte, 12)
	tFile.mtime = make([]byte, 12)
	tFile.chksum = make([]byte, 8)
	//tFile.typeflag=make( []byte,1)
	tFile.linkname = make([]byte, 100)
	tFile.magic = make([]byte, 8)
	tFile.uname = make([]byte, 32)
	tFile.gname = make([]byte, 32)
	tFile.devmajor = make([]byte, 8)
	tFile.devminor = make([]byte, 8)
	tFile.prefix = make([]byte, 167)
	//Find parts of struct given 512 byte block
	copy(tFile.name, b1[:100])
	copy(tFile.mode, b1[100:108])
	copy(tFile.uid, b1[108:116])
	copy(tFile.gid, b1[116:124])
	copy(tFile.size, b1[124:136])
	copy(tFile.mtime, b1[136:148])
	copy(tFile.chksum, b1[148:156])
	tFile.typeflag = b1[156]
	copy(tFile.linkname, b1[157:257])
	copy(tFile.magic, b1[257:265])
	copy(tFile.uname, b1[265:297])
	copy(tFile.gname, b1[297:329])
	copy(tFile.devmajor, b1[329:337])
	copy(tFile.devminor, b1[337:345])
	copy(tFile.prefix, b1[345:])
}

func (tFile *TarFile2) Create(b1 []byte) {

	//Find parts of struct given 512 byte block
	tFile.name = MyBS(b1[:100])
	tFile.mode = MyBS(b1[100:108])
	tFile.uid = MyBtoO(b1[108:116])
	tFile.gid = MyBtoO(b1[116:124])
	tFile.size = MyBtoO(b1[124:136])
	tFile.mtime = MyBtoO(b1[136:148])
	tFile.chksum = MyBtoO(b1[148:156])
	tFile.typeflag = b1[156]
	tFile.linkname = MyBS(b1[157:257])
	tFile.magic = MyBS(b1[257:265])
	tFile.uname = MyBS(b1[265:297])
	tFile.gname = MyBS(b1[297:329])
	tFile.devmajor = MyBtoO(b1[329:337])
	tFile.devminor = MyBtoO(b1[337:345])
	tFile.prefix = MyBS(b1[345:])
}

func PrintPerm(perm int64) {
	if (perm & 4) != 0 {
		fmt.Printf("r")
	} else {
		fmt.Printf("-")
	}
	if (perm & 2) != 0 {
		fmt.Printf("w")
	} else {
		fmt.Printf("-")
	}
	if (perm & 1) != 0 {
		fmt.Printf("x")
	} else {
		fmt.Printf("-")
	}
}

func (tfile TarFile) MyLs() { //do ls -l
	dflag := false
	if tfile.typeflag == '5' { //file is directory
		dflag = true //want to print file blue
		fmt.Printf("d")
	} else {
		fmt.Printf("-")
	}
	mode := string(tfile.mode[len(tfile.mode)-4 : len(tfile.mode)-1])
	uperm, _ := strconv.ParseInt(string(mode[0]), 8, 64)
	gperm, _ := strconv.ParseInt(string(mode[1]), 8, 64)
	operm, _ := strconv.ParseInt(string(mode[2]), 8, 64)
	PrintPerm(uperm)
	PrintPerm(gperm)
	PrintPerm(operm)
	if dflag {
		fmt.Printf(" l ") //file is directory - need number of files in directory
	} else {
		fmt.Printf(" 1 ")
	}
	fmt.Printf("%s %s %d ", string(tfile.uname), string(tfile.gname), MyBtoO(tfile.size))
	//fmt.Printf("%o%o%o", uperm, gperm, operm)

	if dflag {
		fmt.Printf("\x1b[01;34m%s\n\x1b[0m", MyBS(tfile.name)) //print directory blue
	} else {
		fmt.Printf("%s\n", MyBS(tfile.name))
	}
}

func MyBS(bar []byte) string { // convert byte array to string based on null byte
	s := string(bar)
	return s[:strings.Index(s, "\x00")]
}

func MyBtoO(bar []byte) int64 { //convert byte array to octal number
	output, _ := strconv.ParseInt(string(bar)[:len(bar)-1], 8, 64)

	return output
}

func MyBtoI(bar []byte) int64 { //convert byte array to base ten number
	output, _ := strconv.ParseInt(string(bar)[:len(bar)-1], 10, 64)

	return output
}

func main() {
	var tfile TarFile
	var tfile2 TarFile2
	//var sread int64 = 0 //keep track if next block should be read or not
	dat, err := os.Open("../cs671/x.tar")
	if err != nil {
		panic(err)
	}

	b1 := make([]byte, 512)
	_, err2 := dat.Read(b1)
	if err2 != nil {
		panic(err2)
	}

	tfile.Create(b1)
	tfile2.Create(b1)
	fmt.Println(MyBS(tfile.name))
	fmt.Println(tfile2.name)
}
