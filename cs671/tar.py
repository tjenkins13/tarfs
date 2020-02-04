#!/usr/bin/env python3

import sys
import struct


if __name__ == "__main__":

    tarfile = "x.tar" if len(sys.argv) == 1 else sys.argv[1]

    with open(tarfile, "rb") as fd:
        while True:
            fname = fd.read(100).partition(b'\0')[0].decode()
            if len(fname) == 0:
                break
            modes = fd.read(8)
            print(fname)
            print(modes.decode())

            uid = int(fd.read(8).partition(b'\0')[0].decode(), 8)
            gid = int(fd.read(8).partition(b'\0')[0].decode(), 8)
            print(uid,gid)
            size = int(fd.read(12).partition(b'\0')[0].decode(), 8)
            print(size)
            mtime = int(fd.read(12).partition(b'\0')[0].decode(), 8)
            print(mtime)
            fd.read(8)
            flag = int(fd.read(1))
            print(flag)
            fd.read(100+8+32+32+8+8+167)
            if size == 0:
                continue
            data = fd.read(512)
            print('<----',size,'byte file starts here')
            print(data[:size].decode(),end="<---- ends here\n")

