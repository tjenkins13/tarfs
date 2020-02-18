package tree

import (
	"fmt"
	"github.com/tjenkins13/tarfs/tinfo"
	"strings"
	//"time"
)

type TarFile = tinfo.TarFile

type Node struct {
	Data       TarFile
	parent     *Node
	childdir   []Node
	childfiles []TarFile
}

/*
type Sarray []string

func (s Sarray) contains(v string) bool {
	for _, word := range s {
		if word == v {
			return true
		}
	}
	return false
}
*/

func Create_Root() Node {
	var x TarFile
	var n Node
	x.Name = ""
	x.Mode = "755" //for now, might change it to take permission of tar file later
	x.Size = 0
	//x.Mtime = int64(time.Now())
	x.Typeflag = '5'
	n.Data = x
	n.parent = nil
	return n
}

func Create(data TarFile) Node {
	var n Node
	n.Data = data
	n.parent = nil
	return n
}

func Create2(data TarFile) *Node {
	var n Node
	n.Data = data
	n.parent = nil
	return &n
}

//need work
func (u *Node) Insert(data TarFile) {
	path := strings.SplitN(data.Name, "/", -1)
	currpath := strings.SplitN(u.Data.Name, "/", -1)
	//fmt.Println("->", path[len(path)-2])
	//fmt.Println(currpath[len(path)-2])
	if path[len(path)-1] == "" {
		path = path[:len(path)-1]
	}
	if currpath[len(currpath)-1] == "" {
		currpath = currpath[:len(currpath)-1]
	}
	c := 0
	//fmt.Println("-------> Enterring", data.Name, "from", u.Data.Name)
	//fmt.Println(path, len(path), currpath, len(currpath))
	//fmt.Println(u.childdir)
	if len(path) > len(currpath)+1 { //this means that we will need to go through some directories
		for c < len(currpath) && c < len(path) {
			if currpath[c] != path[c] {
				break
			}
			//fmt.Println(currpath[c] == path[c])
			c++ //c has index of new part of path to explore
		}
		//fmt.Println("-->", u.Data.Name, data.Name, c)
		nxtpath := strings.Join(path[:c+1], "/") + "/"

		//fmt.Println("----------", nxtpath)
		for i, file := range u.childdir {
			if file.Data.Name == nxtpath {
				//fmt.Println()
				//fmt.Println("Found next spot to search")
				//fmt.Println(file.Data.Name)
				//fmt.Println()
				u.childdir[i].Insert(data)
				return
			}

		}
		//fmt.Println(path, data.Name, nxtpath)

	} else { //we need to append new directory or file to our tree
		//fmt.Println("-------> Inputting", data.Name, "into", u.Data.Name)

		if data.Typeflag == '5' { //We have another directory
			u.childdir = append(u.childdir, Create(data))
			u.childdir[len(u.childdir)-1].parent = u
		} else { //file is regular--might check for other typeflags
			u.childfiles = append(u.childfiles, data)
			u.Data.Link++
		}
	}
}

func (u Node) Search(name string) (*TarFile, error) {
	path := strings.SplitN(name, "/", -1)
	currpath := strings.SplitN(u.Data.Name, "/", -1)
	//fmt.Println("->", path[len(path)-2])
	//fmt.Println(currpath[len(path)-2])
	if path[len(path)-1] == "" {
		path = path[:len(path)-1]
	}
	if currpath[len(currpath)-1] == "" {
		currpath = currpath[:len(currpath)-1]
	}
	c := 0
	//fmt.Println("-------> Enterring", data.Name, "from", u.Data.Name)
	//fmt.Println(path, len(path), currpath, len(currpath))
	//fmt.Println(u.childdir)
	if len(path) > len(currpath)+1 { //this means that we will need to go through some directories
		for c < len(currpath) && c < len(path) {
			if currpath[c] != path[c] {
				break
			}
			//fmt.Println(currpath[c] == path[c])
			c++ //c has index of new part of path to explore
		}
		//fmt.Println("-->", u.Data.Name, data.Name, c)
		nxtpath := strings.Join(path[:c+1], "/") + "/"

		//fmt.Println("----------", nxtpath)
		for i, file := range u.childdir {
			if file.Data.Name == nxtpath {
				//fmt.Println()
				//fmt.Println("Found next spot to search")
				//fmt.Println(file.Data.Name)
				//fmt.Println()
				return u.childdir[i].Search(name)

			}
			//fmt.Println(path, data.Name, nxtpath)
		}
		return nil, fmt.Errorf("Couldn't find file")

	} else { //we need to append new directory or file to our tree
		//fmt.Println("-------> Inputting", data.Name, "into", u.Data.Name)

		for _, file := range u.childdir {
			if file.Data.Name == name {
				return &file.Data, nil
			}
		}
		for _, file := range u.childfiles {
			if file.Name == name {
				return &file, nil
			}
		}
	}
	//	return nil, &errorString{"Couldn't find file"}
	return nil, fmt.Errorf("Couldn't find file")

}

func (u Node) Print() {
	if u.Data.Name != "" {
		fmt.Println("Beginning------->", u.Data.Name)
		u.Data.MyLs()
	}
	for _, file := range u.childfiles {
		//fmt.Println(file.Name)
		file.MyLs()
	}
	fmt.Println("Ending---------->", u.Data.Name)
	for _, file := range u.childdir {
		fmt.Println(file.Data.Name, "has", len(file.childfiles), "files and", len(file.childdir), "subdirectories")
		file.Print()
	}
}
