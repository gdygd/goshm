package shmlinux

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// System call constants.
const SYS_SHMAT = 30
const SYS_SHMCTL = 31
const SYS_SHMDT = 67
const SYS_SHMGET = 29
const (
	// sysShmAt  = syscall.SYS_SHMAT
	// sysShmCtl = syscall.SYS_SHMCTL
	// sysShmDt  = syscall.SYS_SHMDT
	// sysShmGet = syscall.SYS_SHMGET
	sysShmCtl = SYS_SHMCTL
	sysShmAt  = SYS_SHMAT
	sysShmDt  = SYS_SHMDT
	sysShmGet = SYS_SHMGET
)

const (
	// Create key if key does not exist.
	IPC_CREAT = 01000

	// Private key.
	IPC_PRIVATE = 0

	MEM_READWRITE = 0666

	// Remove identifier.
	IPC_RMID = 0
	// Set `ipc_perm` options.
	IPC_SET = 1
	// Get `ipc_perm' options.
	IPC_STAT = 2

	USERKEY = 0x1234
)

type Linuxshm struct {
	Id   uintptr
	Addr uintptr // memory pointer
	Size int     // allocate memory size
	Skey int     // shared memory key
}

func NewLinuxShm() *Linuxshm {
	return &Linuxshm{}
}

func (m *Linuxshm) InitShm(skey int, size int) {
	m.Size = size
	m.Skey = skey
}

func (m *Linuxshm) CreateShm() error {

	//fmt.Println("CreateShm ..")

	id, _, errno := syscall.Syscall(sysShmGet, uintptr(int32(m.Skey)), uintptr(int32(m.Size)), uintptr(int32(MEM_READWRITE|IPC_CREAT)))
	//fmt.Println("CreateShm#1 : ", id)
	//if int(id) == -1 {
	if int(id) <= 0 {
		//fmt.Println("CreateShm#2 : ", id)
		// Check shm already was made memory
		id, _, errno = syscall.Syscall(sysShmGet, uintptr(int32(m.Skey)), uintptr(int32(m.Size)), uintptr(int32(MEM_READWRITE)))

		//fmt.Println("CreateShm#3 : ", id)
		if int(id) == -1 {
			errmsg := fmt.Sprintf("CreateShm..err: %s", errno.Error())
			err := os.NewSyscallError(errmsg, nil)
			closehandleL(id)
			return err
		}
		//fmt.Println("CreateShm#4 : ", id)
	}

	//fmt.Println("CreateShm#5 ..id:", id)

	m.Id = id
	//return int(id), nil
	return nil
}

func (m *Linuxshm) AttachShm() error {

	//fmt.Println("AttachShm ..id:", m.Id)

	addr, _, errno := syscall.Syscall(sysShmAt, uintptr(int32(m.Id)), 0, uintptr(int32(MEM_READWRITE)))
	if int(addr) == -1 {
		errmsg := fmt.Sprintf("AttachShm..err: %s", errno.Error())
		err := os.NewSyscallError(errmsg, nil)
		closehandleL(m.Id)
		return err
	}

	//fmt.Println("AttachShm ..addr:", addr)

	m.Addr = addr

	// shared memory pointer conver User define type
	// ex) m.UserData = (*SharedMem)(unsafe.Pointer(addr))
	// ex) UserData => type User struct { Name [10]byte, Age int }

	return nil
}

func (m *Linuxshm) DeleteShm() error {
	// detachshm
	// close handle
	err := m.detachShm()
	return err
}

func (m *Linuxshm) detachShm() error {
	// close handle
	// init Addr	(Addr memory pointer)
	// init SKey	(SKey : shared memory key)

	result, _, errno := syscall.Syscall(sysShmDt, uintptr(unsafe.Pointer(m.Addr)), 0, 0)
	if int(result) == -1 {
		errmsg := fmt.Sprintf("detachShm..err: %s", errno.Error())
		err := os.NewSyscallError(errmsg, nil)
		closehandleL(m.Id)
		return err
	}

	m.Addr = 0
	closehandleL(m.Id)

	return nil
}

func closehandleL(id uintptr) (int, error) {

	if id <= 0 {
		return 0, nil
	}

	result, _, errno := syscall.Syscall(sysShmCtl, uintptr(int32(id)), uintptr(int32(IPC_RMID)), uintptr(unsafe.Pointer(nil)))
	if int(result) == -1 {
		return -1, errno
	}

	return int(result), nil

}
