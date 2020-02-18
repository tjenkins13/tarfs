package tinfo

import "fmt"

//import "bufio"
//import "io"
//import "io/ioutil"
//import "os"
import "strconv"
import "strings"

//import "sort"

import "time"

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

//var MAX_OPEN_FILE int = 10

type TarFile struct {
	Name     string
	Mode     string
	Uid      int64
	Gid      int64
	Size     int64
	Mtime    int64
	Chksum   int64
	Typeflag byte
	Linkname string
	Magic    string
	Uname    string
	Gname    string
	Devmajor int64
	Devminor int64
	Prefix   string
	Link     int64 //I added to keep track of num links for ls
	Blockno  int64 //What block is tarfile header in
	Offset   int64 //where are we within the file
}

type ByNameLen []TarFile

func (a ByNameLen) Len() int           { return len(a) }
func (a ByNameLen) Less(i, j int) bool { return len(a[i].Name) < len(a[j].Name) }
func (a ByNameLen) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//var Fd_arr []TarFile = make(TarFile, 0, MAX_OPEN_FILE) //want to be lower case

func (tFile *TarFile) Create(b1 []byte, blockno int64) {

	//Find parts of struct given 512 byte block
	tFile.Name = MyBS(b1[:100])
	tFile.Mode = MyBS(b1[100:108])
	tFile.Mode = tFile.Mode[len(tFile.Mode)-3 : len(tFile.Mode)]
	tFile.Uid = MyBtoO(b1[108:116])
	tFile.Gid = MyBtoO(b1[116:124])
	tFile.Size = MyBtoO(b1[124:136])
	tFile.Mtime = MyBtoO(b1[136:148])
	tFile.Chksum = MyBtoO(b1[148:156])
	tFile.Typeflag = b1[156]
	tFile.Linkname = MyBS(b1[157:257])
	tFile.Magic = MyBS(b1[257:265])
	tFile.Uname = MyBS(b1[265:297])
	tFile.Gname = MyBS(b1[297:329])
	tFile.Devmajor = MyBtoO(b1[329:337])
	tFile.Devminor = MyBtoO(b1[337:345])
	tFile.Prefix = MyBS(b1[345:])
	tFile.Link = 1
	tFile.Blockno = blockno
	tFile.Offset = 0
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
	if tfile.Typeflag == '5' { //file is directory
		dflag = true //want to print file blue
		fmt.Printf("d")
	} else {
		fmt.Printf("-")
	}
	uperm, _ := strconv.ParseInt(string(tfile.Mode[0]), 8, 64)
	gperm, _ := strconv.ParseInt(string(tfile.Mode[1]), 8, 64)
	operm, _ := strconv.ParseInt(string(tfile.Mode[2]), 8, 64)
	PrintPerm(uperm)
	PrintPerm(gperm)
	PrintPerm(operm)
	//if dflag {
	fmt.Printf(" %d ", tfile.Link) //file is directory - need number of files in directory
	//} else {
	//	fmt.Printf(" 1 ")
	//}
	fmt.Printf("%s %s %d ", tfile.Uname, tfile.Gname, tfile.Size)
	//fmt.Printf("%o%o%o", uperm, gperm, operm)
	t := time.Unix(tfile.Mtime, 0)
	fmt.Printf("%s %d %d ", t.Month(), t.Day(), t.Year())
	path := strings.SplitN(tfile.Name, "/", -1)
	if path[len(path)-1] == "" {
		path = path[:len(path)-1]
	}
	if dflag {
		//fmt.Printf("\x1b[01;34m%s\n\x1b[0m", tfile.Name) //print directory blue
		fmt.Printf("\x1b[01;34m%s\x1b[0m/\n", path[len(path)-1]) //print directory blue
	} else {
		//fmt.Printf("%s\n", tfile.Name)
		fmt.Printf("%s\n", path[len(path)-1])
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
