// Copyright (C) 2013-2018 by Maxim Bublis <b@codemonkey.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package uuid

import (
	"unsafe"
)

// byteSliceToString uses pointer manipulation tricks to take a slice of bytes
// and convert it into a string. It should be used to avoid copying the slice.
// It should also cause the compiler's escape analysis to keep the slice on
// the stack if it would otherwise be stack allocated.
func byteSliceToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// stringToByteSlice uses pointer manipulation tricks to take a string and
// return the underlying slice of bytes referenced by the string. This is
// useful when wanting to avoid extra allocations.
func stringToByteSlice(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
