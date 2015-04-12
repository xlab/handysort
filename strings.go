// Copyright 2014 Maxim Kouprianov. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

/*
Package handysort implements an alphanumeric string comparison function
in order to sort alphanumeric strings correctly.

Default sort (incorrect):
	abc1
	abc10
	abc12
	abc2

Handysort:
	abc1
	abc2
	abc10
	abc12

Please note, that handysort is about 5x-8x times slower
than a simple sort, so use it wisely.
*/
package handysort

import "unicode/utf8"

// Strings implements the sort interface, sorts an array
// of the alphanumeric strings in decreasing order.
type Strings []string

func (a Strings) Len() int           { return len(a) }
func (a Strings) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Strings) Less(i, j int) bool { return StringLess(a[i], a[j]) }

// StringLess compares two alphanumeric strings correctly.
func StringLess(s1, s2 string) (less bool) {
	n1, n2 := make([]rune, 0, 20), make([]rune, 0, 20)
	var n1NonEmpty, n2NonEmpty bool

	for i, j := 0, 0; i < len(s1) || j < len(s2); {
		var r1, r2 rune
		var w1, w2 int
		var d1, d2 bool

		// read rune from former string available
		if i < len(s1) {
			r1, w1 = utf8.DecodeRuneInString(s1[i:])
			i += w1

			// if digit, accumulate
			if d1 = ('0' <= r1 && r1 <= '9'); d1 {
				n1 = append(n1, r1)
				n1NonEmpty = true
			}
		}

		// read rune from latter string if available
		if j < len(s2) {
			r2, w2 = utf8.DecodeRuneInString(s2[j:])
			j += w2

			// if digit, accumulate
			if d2 = ('0' <= r2 && r2 <= '9'); d2 {
				n2 = append(n2, r2)
				n2NonEmpty = true
			}
		}

		// if have rune and other non-digit rune
		if (!d1 || !d2) && r1 > 0 && r2 > 0 {
			if n1NonEmpty && n2NonEmpty {
				// compare digits in accumulators
				return compareByDigits(n1, n2)
			}

			// detect if non-digit rune from former or latter
			if r1 != r2 {
				return r1 < r2
			}
		}
	}

	if n1NonEmpty || n2NonEmpty {
		// reached both strings ends, compare numeric accumulators
		return compareByDigits(n1, n2)
	}

	// last hope
	return len(s1) < len(s2)
}

// Compare two numeric fields by their digits
func compareByDigits(n1, n2 []rune) (less bool) {
	offset := len(n2) - len(n1)
	n1n2 := offset < 0 // len(n1) > len(n2)
	if n1n2 {
		// if n1 longer, swap with n2
		offset = -offset
		n1, n2 = n2, n1
	}

	var j int
	// len(n1) always be <= len(n2)
	for i := range n2 {
		var r1 rune
		if offset == 0 {
			// begin actual read
			r1 = n1[j]
			j++
		} else {
			// emulate zero-padding
			r1 = '0'
			offset--
		}

		r2 := n2[i]
		if r1 != r2 {
			if n1n2 {
				return r2 < r1 // actually r1 < r2
			}
			return r1 < r2
		}
	}

	// numeric value equals, compare by length
	if n1n2 {
		// n1 was > n2
		return false
	}
	// eval a comparison only if n1 known to be <= n2
	return len(n1) < len(n2)
}
