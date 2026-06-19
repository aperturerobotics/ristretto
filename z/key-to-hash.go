//go:build !js

/*
 * SPDX-FileCopyrightText: © 2017-2025 Istari Digital, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package z

import (
	"reflect"

	"github.com/cespare/xxhash/v2"
)

// KeyToHash hashes a cache key into the two uint64 hashes ristretto stores.
// Exact builtin key types take the type switch; a named key type whose
// underlying type satisfies Key is resolved with reflect so callers may use
// types such as `type MyKey string`. The js build replaces this with a
// reflect-free variant.
//
// TODO: Figure out a way to re-use memhash for the second uint64 hash, we
// already know that appending bytes isn't reliable for generating a second hash
// (see Ristretto PR #88). We also know that while the Go runtime has a runtime
// memhash128 function, it's not possible to use it to generate [2]uint64 or
// anything resembling a 128bit hash, even though that's exactly what we need in
// this situation.
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
		// Handle custom types with underlying types (e.g., type MyKey string).
		v := reflect.ValueOf(key)
		switch v.Kind() {
		case reflect.Uint64:
			return v.Uint(), 0
		case reflect.String:
			s := v.String()
			return MemHashString(s), xxhash.Sum64String(s)
		case reflect.Slice:
			if v.Type().Elem().Kind() == reflect.Uint8 {
				b := v.Bytes()
				return MemHash(b), xxhash.Sum64(b)
			}
		case reflect.Uint8:
			return v.Uint(), 0
		case reflect.Uint:
			return v.Uint(), 0
		case reflect.Int:
			return uint64(v.Int()), 0
		case reflect.Int32:
			return uint64(v.Int()), 0
		case reflect.Uint32:
			return v.Uint(), 0
		case reflect.Int64:
			return uint64(v.Int()), 0
		}
		panic("Key type not supported")
	}
}
