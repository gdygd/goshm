# 'Go Shared memory library'

`goshm` is a Shared memory library for linux and windows

## Installation

    go get github.com/gdygd/goshm

## example
 1. Create Shared Memory
 2. Example of reading and writing to a shared memory using a User define structure 
   - "type ShmMem struct { ... }"

### process1 example

  ```go
package main

import (
	"fmt"
	"unsafe"

	"github.com/gdygd/goshm/shmlinux"
    //"github.com/gdygd/goshm/shmwin"   // windows lib package
)

type ShmMem struct {
	Name [10]byte
	Num1 int
	Num2 int
	Addr [10]byte
}

var SharedMem *ShmMem

const skey = 0x1234 // shared memory key

func main() {
    // linux os
	shminst := shmlinux.NewLinuxShm()   
	shminst.InitShm(skey, 1024)

    // windows
    // shminst := shmwin.NewWinShm()
	// shminst.InitShm("shmygd", 1024)

    // The code below is the same for all Windows Linux.
    //--------------------------------------------------------------
    // shminst.CreateShm() // Create or open shared memory
    // shminst.AttachShm() // Attach shared memory
    // shminst.DeleteShm() // Detatch and Delete Shared memory
    //--------------------------------------------------------------

	err := shminst.CreateShm()
	if err != nil {
		fmt.Println("CreateShm err : ", err)
	}

	err = shminst.AttachShm()
	if err != nil {
		fmt.Println("AttachShm err : ", err)
	}
	SharedMem = (*ShmMem)(unsafe.Pointer(shminst.Addr))

	copy(SharedMem.Addr[:], []byte(string("1234aa")))
	copy(SharedMem.Name[:], []byte(string("YGD")))
	SharedMem.Num1 = 11
	SharedMem.Num2 = 99

	var n int
	fmt.Scanf("%d", &n)

	fmt.Printf("\nmem#1 : %v\n", SharedMem)

	err = shminst.DeleteShm()
	if err != nil {
		fmt.Println("DeleteShm err:", err)
	}
}

  ```

### process2 example

  ```go
package main

import (
	"fmt"
	"unsafe"

	"github.com/gdygd/goshm/shmlinux"
    //"github.com/gdygd/goshm/shmwin"   // windows lib package
)

type ShmMem struct {
	Name [10]byte
	Num1 int
	Num2 int
	Addr [10]byte
}

var SharedMem *ShmMem

const skey = 0x1234

func main() {
	shminst := shmlinux.NewLinuxShm()

	shminst.InitShm(skey, 1024)

    // windows
    // shminst := shmwin.NewWinShm()
	// shminst.InitShm("shmygd", 1024)

    // The code below is the same for all Windows Linux.
    //--------------------------------------------------------------
    // shminst.CreateShm() // Create or open shared memory
    // shminst.AttachShm() // Attach shared memory
    // shminst.DeleteShm() // Detatch and Delete Shared memory
    //--------------------------------------------------------------

	err := shminst.CreateShm()
	if err != nil {
		fmt.Println("CreateShm err : ", err)
	}

	err = shminst.AttachShm()
	if err != nil {
		fmt.Println("AttachShm err : ", err)
	}
	SharedMem = (*ShmMem)(unsafe.Pointer(shminst.Addr))

	fmt.Printf("mem#2 : %v\n", SharedMem)

	SharedMem.Num1 = 1
	SharedMem.Num2 = 2

	err = shminst.DeleteShm()
	if err != nil {
		fmt.Println("DeleteShm err:", err)
	}
}


  ```