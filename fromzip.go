//I got this from https://github.com/bazil/zipfs/blob/master/main.go#L78
//and I am editing it to work for tar files using what I wrote

package main

import (
	//	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"

	"github.com/tjenkins13/tarfs/tinfo"
	"github.com/tjenkins13/tarfs/tree"
)

//Stuff I put in
type TarFile = tinfo.TarFile
type Node = tree.Node

type OpenStruct struct {
	tname    string //absolute path to tar file
	mloc     string //absolute path to where you are mounting tar file
	fname    string //name of file within tar file you want
	tree     *Node  //head of tree for filesystem
	allfiles []TarFile
	fd       *os.File //file descriptor for tar file
}

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

//End of stuff I put in
// We assume the zip file contains entries for directories too.

var progName = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", progName)
	fmt.Fprintf(os.Stderr, "  %s TAR MOUNTPOINT\n", progName)
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(progName + ": ")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		os.Exit(2)
	}
	path := flag.Arg(0)
	mountpoint := flag.Arg(1)
	if err := mount(path, mountpoint); err != nil {
		log.Fatal(err)
	}
}

func (a *OpenStruct) Close() error {
	//a.fd.Close()
	return nil
}

func OpenTar(filename string) (*OpenStruct, error) {
	var tfile []TarFile
	var t TarFile
	var block int64 = 0
	var head Node
	var a OpenStruct
	a.tname = filename
	sread := 0
	dat, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	a.fd = dat
	b1 := make([]byte, 512)
	_, err2 := dat.Read(b1)
	if err2 != nil {
		return nil, err2
	}
	//fmt.Println(string(b1[:100]))
	//fmt.Println(b1)
	for !allZero(b1) {
		if sread == 0 {
			t.Create(b1, block)
			//fmt.Println(t.Name)
			tfile = append(tfile, t)

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
		block++
	}
	sort.Sort(ByNameLen(tfile)) // sort by length of file name for insertion into tree
	//get tree of filesystem and put it in array
	head = tree.Create_Root()
	//head = tree.Create(tfile[0])
	//fmt.Println(head.Data.Name)
	//for _, file := range tfile[1:] {
	for _, file := range tfile {
		head.Insert(file)
		//fmt.Println(file.Name)
	}
	a.allfiles = tfile
	a.tree = &head

	//head.Print()
	return &a, nil

}

func mount(path, mountpoint string) error {
	//	archive, err := zip.OpenReader(path)
	archive, err := OpenTar(path)
	if err != nil {
		return err
	}
	defer archive.Close()

	c, err := fuse.Mount(mountpoint)
	if err != nil {
		return err
	}
	defer c.Close()

	filesys := &FS{
		archive: archive,
	}
	if err := fs.Serve(c, filesys); err != nil {
		return err
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		return err
	}

	return nil
}

type FS struct {
	//archive *zip.Reader
	archive *OpenStruct
}

var _ fs.FS = (*FS)(nil)

func (f *FS) Root() (fs.Node, error) {
	n := &Dir{
		archive: f.archive,
	}
	return n, nil
}

type Dir struct {
	//archive *zip.Reader
	archive *OpenStruct
	// nil for the root directory, which has no entry in the zip
	file *TarFile
}

var _ fs.Node = (*Dir)(nil)

//func zipAttr(f *zip.File, a *fuse.Attr) {
func zipAttr(f *TarFile, a *fuse.Attr) {
	a.Size = uint64(f.Size)
	mode, _ := strconv.ParseInt(f.Mode, 8, 64)
	if f.Typeflag == '5' {
		mode |= int64(uint32(os.ModeDir))
	}
	a.Mode = os.FileMode(uint32(mode))
	a.Mtime = time.Unix(f.Mtime, 0)
	a.Ctime = time.Unix(f.Mtime, 0)
	a.Crtime = time.Unix(f.Mtime, 0)
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	if d.file == nil {
		// root directory
		a.Mode = os.ModeDir | 0755
		return nil
	}
	zipAttr(d.file, a)
	return nil
}

var _ = fs.NodeRequestLookuper(&Dir{})

//Need to work on this
func (d *Dir) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (fs.Node, error) {
	path := req.Name
	if d.file != nil {
		path = d.file.Name + path
	}
	for _, f := range d.archive.allfiles {
		switch {
		case f.Name == path: //File is normal file
			child := &File{
				file:    &f,
				archive: d.archive,
			}
			return child, nil
		case f.Name[:len(f.Name)-1] == path && f.Name[len(f.Name)-1] == '/': //File is directory
			child := &Dir{
				archive: d.archive,
				file:    &f,
			}
			return child, nil
		}
	}
	return nil, fuse.ENOENT
}

var _ = fs.HandleReadDirAller(&Dir{})

//Need to work on this
func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	prefix := ""
	if d.file != nil {
		prefix = d.file.Name
	}

	var res []fuse.Dirent
	for _, f := range d.archive.allfiles {
		if !strings.HasPrefix(f.Name, prefix) {
			continue
		}
		name := f.Name[len(prefix):]
		if name == "" {
			// the dir itself, not a child
			continue
		}
		if strings.ContainsRune(name[:len(name)-1], '/') {
			// contains slash in the middle -> is in a deeper subdir
			continue
		}
		var de fuse.Dirent
		if name[len(name)-1] == '/' {
			// directory
			name = name[:len(name)-1]
			de.Type = fuse.DT_Dir
		}
		de.Name = name
		res = append(res, de)
	}
	return res, nil
}

type File struct {
	//	file *zip.File
	file    *TarFile
	archive *OpenStruct
}

var _ fs.Node = (*File)(nil)

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	zipAttr(f.file, a)
	return nil
}

var _ = fs.NodeOpener(&File{})

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	/*r, err := f.file.Open()
	if err != nil {
		return nil, err
	}*/
	// individual entries inside a zip file are not seekable
	resp.Flags |= fuse.OpenNonSeekable
	return &FileHandle{r: f.archive, f: f.file}, nil
}

type FileHandle struct {
	//	r io.ReadCloser
	r *OpenStruct
	f *TarFile
}

var _ fs.Handle = (*FileHandle)(nil)

var _ fs.HandleReleaser = (*FileHandle)(nil)

func (fh *FileHandle) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	return fh.r.Close()
}

var _ = fs.HandleReader(&FileHandle{})

func (fh *FileHandle) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	// We don't actually enforce Offset to match where previous read
	// ended. Maybe we should, but that would mean'd we need to track
	// it. The kernel *should* do it for us, based on the
	// fuse.OpenNonSeekable flag.
	//
	// One exception to the above is if we fail to fully populate a
	// page cache page; a read into page cache is always page aligned.
	// Make sure we never serve a partial read, to avoid that.
	buf := make([]byte, req.Size)
	fh.r.fd.Seek((fh.f.Blockno+1)*512+fh.f.Offset, 0)
	n, err := fh.r.fd.Read(buf)
	fh.r.fd.Seek(0, 0)
	if err == io.ErrUnexpectedEOF || err == io.EOF {
		err = nil
	}
	resp.Data = buf[:n]
	return err
}
