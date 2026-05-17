//go:build tinygo && !js && !wasip1

/*
 * SPDX-FileCopyrightText: © 2017-2025 Istari Digital, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package z

import (
	"errors"
	"os"
)

var errMmapUnsupported = errors.New("mmap unsupported on tinygo")

func mmap(fd *os.File, writeable bool, size int64) ([]byte, error) {
	return nil, errMmapUnsupported
}

func munmap(b []byte) error {
	return errMmapUnsupported
}

func madvise(b []byte, readahead bool) error {
	return errMmapUnsupported
}

func msync(b []byte) error {
	return errMmapUnsupported
}
