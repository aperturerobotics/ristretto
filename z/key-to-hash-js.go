//go:build js

/*
 * SPDX-FileCopyrightText: © 2017-2025 Istari Digital, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package z

import "github.com/cespare/xxhash/v2"

// KeyToHash hashes a cache key without reflect for the js browser build. The Key
// constraint guarantees a builtin underlying type, and the type switch covers
// every exact builtin key type. A named key type whose dynamic type is not a
// builtin reaches the default and is unsupported here: in the browser closure
// ristretto keys are always exact builtins (for example the string block cache
// key), so the reflect-based named-type path the native build provides is dead
// weight that would pull reflect into the GoScript bundle.
func KeyToHash[K Key](key K) (uint64, uint64) {
	switch k := any(key).(type) {
	case uint64:
		return k, 0
	case string:
		return MemHashString(k), xxhash.Sum64String(k)
	case []byte:
		return MemHash(k), xxhash.Sum64(k)
	case byte:
		return uint64(k), 0
	case uint:
		return uint64(k), 0
	case int:
		return uint64(k), 0
	case int32:
		return uint64(k), 0
	case uint32:
		return uint64(k), 0
	case int64:
		return uint64(k), 0
	default:
		panic("z: KeyToHash of a named key type is not supported in the js build")
	}
}
