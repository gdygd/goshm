package shmwin

import (
	"os"

	"golang.org/x/sys/windows"
)

type Winshm struct {
	Hnd  windows.Handle
	Addr uintptr // memory pointer
	Size int     // allocate memory size
	Skey string  // shared memory key
}

func NewWinShm() *Winshm {
	return &Winshm{}
}

func (m *Winshm) InitShm(skey string, size int) {
	m.Size = size
	m.Skey = skey
}

func (m *Winshm) CreateShm() error {
	prot := windows.PAGE_READWRITE
	skeyptr, err := windows.UTF16PtrFromString(m.Skey)

	if err != nil {
		return os.NewSyscallError("CreateShm UTF16PtrFromString err:", err)
	}

	h, errno := windows.CreateFileMapping(windows.InvalidHandle, nil, uint32(prot), 0, uint32(m.Size), skeyptr)
	if h == 0 {
		err = os.NewSyscallError("CreateShm..err: ", errno)
		closehandleW(m.Hnd)
		return err
	}

	m.Hnd = h
	return nil
}

func (m *Winshm) AttachShm() error {
	access := uint32(windows.FILE_MAP_READ | windows.FILE_MAP_WRITE)

	addr, errno := windows.MapViewOfFile(m.Hnd, access, 0, 0, uintptr(m.Size))
	if addr == 0 {
		closehandleW(m.Hnd)
		return os.NewSyscallError("MapViewOfFile err:", errno)
	}

	m.Addr = addr // shared memory pointer

	// shared memory pointer conver User define type
	// ex) m.UserData = (*SharedMem)(unsafe.Pointer(addr))
	// ex) UserData => type User struct { Name [10]byte, Age int }

	return nil
}

func (m *Winshm) DeleteShm() error {
	// detachshm
	// close handle
	err := m.detachShm()
	return err
}

func (m *Winshm) detachShm() error {
	// close handle
	// init Addr	(Addr memory pointer)
	// init SKey	(SKey : shared memory key)

	if m.Hnd <= 0 {
		return os.NewSyscallError("DeleteShm fail.. Invalid handle", nil)
	}

	if m.Addr <= 0 {
		return os.NewSyscallError("DeleteShm fail.. Invalid memory address", nil)
	}

	err := windows.UnmapViewOfFile(m.Addr)
	if err != nil {
		err = os.NewSyscallError("UnmapViewOfFile err:", err)
		return err
	}

	m.Addr = 0
	closehandleW(m.Hnd)

	return nil
}

func closehandleW(hnd windows.Handle) {
	if hnd > 0 {
		windows.CloseHandle(hnd)
	}
}
