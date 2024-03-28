package main

import (
	"fmt"
	"unsafe"

	"github.com/gdygd/goshm/shmwin"
)

//var shminst shm.ShmHnd

type ShmMem struct {
	Name [10]byte
	Num1 int
	Num2 int
	Addr [10]byte
}

var SharedMem *ShmMem

func main() {
	shminst := shmwin.NewWinShm()

	shminst.InitShm("shmygd", 1024)

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
