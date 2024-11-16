package util

import (
	"runtime"
	"syscall"
	"unsafe"
)

func getProcessEntry(pid int) (*syscall.ProcessEntry32, error) {
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}

	defer syscall.CloseHandle(snapshot)
	var procEntry syscall.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = syscall.Process32First(snapshot, &procEntry); err != nil {
		return nil, err
	}

	for {
		if procEntry.ProcessID == uint32(pid) {
			return &procEntry, nil
		}

		err = syscall.Process32Next(snapshot, &procEntry)
		if err != nil {
			return nil, err
		}
	}
}

var knownCallers = [2]string{"explorer.exe", "PowerToys.PowerLauncher.exe"}

func IsRunningFromExplorer() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	pe, err := getProcessEntry(syscall.Getppid())
	if err != nil {
		return false
	}

	caller := syscall.UTF16ToString(pe.ExeFile[:])
	for _, exe := range &knownCallers {
		if exe == caller {
			return true
		}
	}

	return false
}
