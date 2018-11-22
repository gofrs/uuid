// Copyright (c) 2018 Andrei Tudor Călin <mail@acln.ro>
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

import v3 "github.com/gofrs/uuid/v3"

// Size of a UUID in bytes.
const Size = v3.Size

// UUID is an array type to represent the value of a UUID, as defined in RFC-4122.
type UUID = v3.UUID

// UUID versions.
const (
	V1 = v3.V1 // Version 1 (date-time and MAC address)
	V2 = v3.V2 // Version 2 (date-time and MAC address, DCE security version)
	V3 = v3.V3 // Version 3 (namespace name-based)
	V4 = v3.V4 // Version 4 (random)
	V5 = v3.V5 // Version 5 (namespace name-based)
)

// UUID layout variants.
const (
	VariantNCS       = v3.VariantNCS
	VariantRFC4122   = v3.VariantRFC4122
	VariantMicrosoft = v3.VariantMicrosoft
	VariantFuture    = v3.VariantFuture
)

// UUID DCE domains.
const (
	DomainPerson = v3.DomainPerson
	DomainGroup  = v3.DomainGroup
	DomainOrg    = v3.DomainOrg
)

// Timestamp is the count of 100-nanosecond intervals since 00:00:00.00,
// 15 October 1582 within a V1 UUID. This type has no meaning for V2-V5
// UUIDs since they don't have an embedded timestamp.
type Timestamp = v3.Timestamp

// TimestampFromV1 returns the Timestamp embedded within a V1 UUID.
// Returns an error if the UUID is any version other than 1.
func TimestampFromV1(u UUID) (Timestamp, error) {
	return v3.TimestampFromV1(u)
}

// Nil is the nil UUID, as specified in RFC-4122, that has all 128 bits set to
// zero.
var Nil = v3.Nil

// Predefined namespace UUIDs.
var (
	NamespaceDNS  = v3.NamespaceDNS
	NamespaceURL  = v3.NamespaceURL
	NamespaceOID  = v3.NamespaceOID
	NamespaceX500 = v3.NamespaceX500
)

// Must is a helper that wraps a call to a function returning (UUID, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//  var packageUUID = v3.Must(v3.FromString("123e4567-e89b-12d3-a456-426655440000"))
func Must(u UUID, err error) UUID {
	return v3.Must(u, err)
}

// NullUUID can be used with the standard sql package to represent a
// UUID value that can be NULL in the database.
type NullUUID = v3.NullUUID

// HWAddrFunc is the function type used to provide hardware (MAC) addresses.
type HWAddrFunc = v3.HWAddrFunc

// Generator provides an interface for generating UUIDs.
type Generator = v3.Generator

// DefaultGenerator is the default UUID Generator used by this package.
var DefaultGenerator = v3.DefaultGenerator

// NewV1 returns a UUID based on the current timestamp and MAC address.
func NewV1() (UUID, error) {
	return v3.NewV1()
}

// NewV2 returns a DCE Security UUID based on the POSIX UID/GID.
func NewV2(domain byte) (UUID, error) {
	return v3.NewV2(domain)
}

// NewV3 returns a UUID based on the MD5 hash of the namespace UUID and name.
func NewV3(ns UUID, name string) UUID {
	return v3.NewV3(ns, name)
}

// NewV4 returns a randomly generated UUID.
func NewV4() (UUID, error) {
	return v3.NewV4()
}

// NewV5 returns a UUID based on SHA-1 hash of the namespace UUID and name.
func NewV5(ns UUID, name string) UUID {
	return v3.NewV5(ns, name)
}

// Gen is a reference UUID generator based on the specifications laid out in
// RFC-4122 and DCE 1.1: Authentication and Security Services. This type
// satisfies the Generator interface as defined in this package.
//
// For consumers who are generating V1 UUIDs, but don't want to expose the MAC
// address of the node generating the UUIDs, the NewGenWithHWAF() function has been
// provided as a convenience. See the function's documentation for more info.
//
// The authors of this package do not feel that the majority of users will need
// to obfuscate their MAC address, and so we recommend using NewGen() to create
// a new generator.
type Gen = v3.Gen

// NewGen returns a new instance of Gen with some default values set. Most
// people should use this.
func NewGen() *Gen {
	return v3.NewGen()
}

// NewGenWithHWAF builds a new UUID generator with the HWAddrFunc provided. Most
// consumers should use NewGen() instead.
//
// This is used so that consumers can generate their own MAC addresses, for use
// in the generated UUIDs, if there is some concern about exposing the physical
// address of the machine generating the UUID.
//
// The Gen generator will only invoke the HWAddrFunc once, and cache that MAC
// address for all the future UUIDs generated by it. If you'd like to switch the
// MAC address being used, you'll need to create a new generator using this
// function.
func NewGenWithHWAF(hwaf HWAddrFunc) *Gen {
	return v3.NewGenWithHWAF(hwaf)
}

// FromBytes returns a UUID generated from the raw byte slice input.
// It will return an error if the slice isn't 16 bytes long.
func FromBytes(input []byte) (UUID, error) {
	return v3.FromBytes(input)
}

// FromBytesOrNil returns a UUID generated from the raw byte slice input.
// Same behavior as FromBytes(), but returns v3.Nil instead of an error.
func FromBytesOrNil(input []byte) UUID {
	return v3.FromBytesOrNil(input)
}

// FromString returns a UUID parsed from the input string.
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (UUID, error) {
	return v3.FromString(input)
}

// FromStringOrNil returns a UUID parsed from the input string.
// Same behavior as FromString(), but returns v3.Nil instead of an error.
func FromStringOrNil(input string) UUID {
	return v3.FromStringOrNil(input)
}
