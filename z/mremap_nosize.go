//go:build (arm64 || arm) && linux && !js && !tinygo
// +build arm64 arm
// +build linux
// +build !js
// +build !tinygo

/*
 * SPDX-FileCopyrightText: © 2017-2025 Istari Digital, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package z

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// mremap is a Linux-specific system call to remap pages in memory. This can be used in place of munmap + mmap.
func mremap(data []byte, size int) ([]byte, error) {
	//nolint:lll
	// taken from <https://github.com/torvalds/linux/blob/f8394f232b1eab649ce2df5c5f15b0e528c92091/include/uapi/linux/mman.h#L8>
	const MREMAP_MAYMOVE = 0x1

	ptr := unsafe.Pointer(unsafe.SliceData(data))
	// For ARM64, the second return argument for SYS_MREMAP is inconsistent (prior allocated size) with
	// other architectures, which return the size allocated
	mmapAddr, _, errno := unix.Syscall6(
		unix.SYS_MREMAP,
		uintptr(ptr),
		uintptr(len(data)),
		uintptr(size),
		uintptr(MREMAP_MAYMOVE),
		0,
		0,
	)
	if errno != 0 {
		return nil, errno
	}

	return unsafe.Slice((*byte)(unsafe.Pointer(mmapAddr)), size), nil
}
