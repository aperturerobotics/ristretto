//go:build tinygo

/*
 * SPDX-FileCopyrightText: © 2017-2025 Istari Digital, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package z

import (
	"sync/atomic"
	"time"

	"github.com/cespare/xxhash/v2"
)

const memHashSeed = 0x9e3779b97f4a7c15

// NanoTime returns the current time in nanoseconds.
func NanoTime() int64 {
	return time.Now().UnixNano()
}

// CPUTicks returns a portable TinyGo timing value.
func CPUTicks() int64 {
	return NanoTime()
}

// MemHash returns a portable TinyGo hash for byte keys.
func MemHash(data []byte) uint64 {
	d := xxhash.NewWithSeed(memHashSeed)
	_, _ = d.Write(data)
	return d.Sum64()
}

// MemHashString returns a portable TinyGo hash for string keys.
func MemHashString(str string) uint64 {
	d := xxhash.NewWithSeed(memHashSeed)
	_, _ = d.WriteString(str)
	return d.Sum64()
}

var fastRandState uint32 = 1

// FastRand returns a cheap process-local pseudo-random value.
func FastRand() uint32 {
	x := atomic.AddUint32(&fastRandState, 0x9e3779b9)
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	return x
}

func Memclr(b []byte) {
	clear(b)
}
