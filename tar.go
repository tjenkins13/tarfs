package main

import "fmt"
//import "bufio"
//import "io"
//import "io/ioutil"
import "os"

/*type TarFile struct{
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
type TarFile struct{
name []byte
mode []byte
uid []byte
gid []byte
size []byte
mtime []byte
chksum []byte
typeflag byte
linkname []byte
magic []byte
uname []byte
gname []byte
devmajor []byte
devminor []byte
prefix []byte
}



func main(){
var tFile TarFile 

fmt.Println("Hello")
dat,err :=os.Open("../tar/cs671/x.tar")
if(err != nil){
panic(err)
}


b1 := make([]byte, 512)
_,err2:=dat.Read(b1)
if err2!=nil{
   panic(err2)
}
//fmt.Println(string(b1[:512]))
tFile.name = make([]byte, 100)
ret := copy(tFile.name,b1[:100])
fmt.Println(ret)
fmt.Println(string(tFile.name))

}