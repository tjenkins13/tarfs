package tinfo

import "fmt"

//import "bufio"
//import "io"
//import "io/ioutil"
//import "os"
import "strconv"
import "strings"

//import "sort"

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

//var MAX_OPEN_FILE int = 10

type TarFile struct {
	Name     string
	mode     string
	uid      int64
	gid      int64
	Size     int64
	mtime    int64
	chksum   int64
	Typeflag byte
	linkname string
	magic    string
	uname    string
	gname    string
	devmajor int64
	devminor int64
	prefix   string
	Link     int64 //I added to keep track of num links for ls
	blockno  int64 //What block is tarfile header in
}

type ByNameLen []TarFile

func (a ByNameLen) Len() int           { return len(a) }
func (a ByNameLen) Less(i, j int) bool { return len(a[i].Name) < len(a[j].Name) }
func (a ByNameLen) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//var Fd_arr []TarFile = make(TarFile, 0, MAX_OPEN_FILE) //want to be lower case

func (tFile *TarFile) Create(b1 []byte, blockno int64) {

	//Find parts of struct given 512 byte block
	tFile.Name = MyBS(b1[:100])
	tFile.mode = MyBS(b1[100:108])
	tFile.mode = tFile.mode[len(tFile.mode)-3 : len(tFile.mode)]
	tFile.uid = MyBtoO(b1[108:116])
	tFile.gid = MyBtoO(b1[116:124])
	tFile.Size = MyBtoO(b1[124:136])
	tFile.mtime = MyBtoO(b1[136:148])
	tFile.chksum = MyBtoO(b1[148:156])
	tFile.Typeflag = b1[156]
	tFile.linkname = MyBS(b1[157:257])
	tFile.magic = MyBS(b1[257:265])
	tFile.uname = MyBS(b1[265:297])
	tFile.gname = MyBS(b1[297:329])
	tFile.devmajor = MyBtoO(b1[329:337])
	tFile.devminor = MyBtoO(b1[337:345])
	tFile.prefix = MyBS(b1[345:])
	tFile.Link = 1
	tFile.blockno = blockno
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
	uperm, _ := strconv.ParseInt(string(tfile.mode[0]), 8, 64)
	gperm, _ := strconv.ParseInt(string(tfile.mode[1]), 8, 64)
	operm, _ := strconv.ParseInt(string(tfile.mode[2]), 8, 64)
	PrintPerm(uperm)
	PrintPerm(gperm)
	PrintPerm(operm)
	//if dflag {
	fmt.Printf(" %d ", tfile.Link) //file is directory - need number of files in directory
	//} else {
	//	fmt.Printf(" 1 ")
	//}
	fmt.Printf("%s %s %d ", tfile.uname, tfile.gname, tfile.Size)
	//fmt.Printf("%o%o%o", uperm, gperm, operm)

	if dflag {
		fmt.Printf("\x1b[01;34m%s\n\x1b[0m", tfile.Name) //print directory blue
	} else {
		fmt.Printf("%s\n", tfile.Name)
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

/*
func allZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

func OpenTar(filename string) int {
	var tfile []TarFile
	var t TarFile
	var block int64 = 0
	sread := 0
	dat, err := os.Open("../cs671/x.tar")
	if err != nil {
		panic(err)
	}

	b1 := make([]byte, 512)
	_, err2 := dat.Read(b1)
	if err2 != nil {
		panic(err2)
	}

	for allZero(b1) {
		if sread == 0 {
			t.Create(b1, block)
			tfile = append(tfile, t)
			block++
			if tfile[len(tfile)-1].Typeflag != '5' { //if it's not a directory we're about to read some data
				sread = int((tfile[len(tfile)-1].Size-1)/512) + 1
			}
		} else {
			sread--
		}
		dat.Read(b1)
	}
	sort.Sort(ByNameLen(tfile)) // sort by length of file name for insertion into tree

	return len(tfile) //for now, will be change to return element of array

}

func CloseTar(fd int) {

}
*/
/*
func main() {
	//	var tfile TarFile
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

	tfile.Create(b1,0)
	fmt.Println(len(tfile.name))
	fmt.Println(tfile.mode[0])
	fmt.Println(tfile.mode[1])
	fmt.Println(tfile.mode[2])
	fmt.Println(tfile.Size)
	fmt.Println(string(tfile.typeflag))
	tfile.MyLs()

	fmt.Println(OpenTar("../cs671/x.tar"))
}
*/
