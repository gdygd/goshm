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
)

type Linuxshm struct {
	Id   int
	Addr uintptr // memory pointer
	Size int     // allocate memory size
	Skey string  // shared memory key
}

func NewLinuxShm() *Linuxshm {
	return &Linuxshm{}
}

func (m *Linuxshm) InitShm(skey string, size int) {
	m.Size = size
	m.Skey = skey
}

func (m *Linuxshm) CreateShm() error {

	id, _, errno := syscall.SyscallN(sysShmGet, uintptr(int32(IPC_CREAT)), uintptr(int32(m.Size)), uintptr(int32(MEM_READWRITE)))
	if int(id) == -1 {
		// Check shm already was made memory
		id, _, errno = syscall.SyscallN(sysShmGet, uintptr(int32(IPC_PRIVATE)), uintptr(int32(m.Size)), uintptr(int32(MEM_READWRITE)))

		if int(id) == -1 {
			errmsg := fmt.Sprintf("CreateShm..err: %s", errno.Error())
			err := os.NewSyscallError(errmsg, nil)
			closehandleL(int(id))
			return err
		}
	}

	//return int(id), nil
	return nil
}

func (m *Linuxshm) AttachShm() error {

	addr, _, errno := syscall.SyscallN(sysShmAt, uintptr(int32(m.Id)), 0, uintptr(int32(MEM_READWRITE)))
	if int(addr) == -1 {
		errmsg := fmt.Sprintf("AttachShm..err: %s", errno.Error())
		err := os.NewSyscallError(errmsg, nil)
		closehandleL(int(m.Id))
		return err
	}

	m.Addr = addr

	// shared memory pointer conver User define type
	// ex) m.UserData = (*SharedMem)(unsafe.Pointer(addr))
	// ex) UserData => type User struct { Name [10]byte, Age int }

	return nil
}

func (m *Linuxshm) DeleteShm() {
	// detachshm
	// close handle
	m.detachShm()
}

func (m *Linuxshm) detachShm() error {
	// close handle
	// init Addr	(Addr memory pointer)
	// init SKey	(SKey : shared memory key)

	result, _, errno := syscall.SyscallN(sysShmDt, uintptr(unsafe.Pointer(m.Addr)), 0, 0)
	if int(result) == -1 {
		errmsg := fmt.Sprintf("detachShm..err: %s", errno.Error())
		err := os.NewSyscallError(errmsg, nil)
		closehandleL(int(m.Id))
		return err
	}

	m.Addr = 0
	closehandleL(m.Id)

	return nil
}

func closehandleL(id int) (int, error) {

	if id <= 0 {
		return 0, nil
	}

	result, _, errno := syscall.SyscallN(sysShmCtl, uintptr(int32(id)), uintptr(int32(IPC_RMID)), uintptr(unsafe.Pointer(nil)))
	if int(result) == -1 {
		return -1, errno
	}

	return int(result), nil

}
