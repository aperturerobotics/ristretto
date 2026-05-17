//go:build !tinygo

/*
 * SPDX-FileCopyrightText: © 2017-2025 Istari Digital, Inc.
 * SPDX-License-Identifier: Apache-2.0
 */

package z

import "os"

func getPageSize() int {
	return os.Getpagesize()
}
