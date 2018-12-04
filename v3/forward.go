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

package v3

import (
	"github.com/gofrs/uuid"
)

// Size of a UUID in bytes.
const Size = uuid.Size

// UUID is an array type to represent the value of a UUID, as defined in RFC-4122.
type UUID = uuid.UUID

// UUID versions.
const (
	V1 = uuid.V1 // Version 1 (date-time and MAC address)
	V2 = uuid.V2 // Version 2 (date-time and MAC address, DCE security version)
	V3 = uuid.V3 // Version 3 (namespace name-based)
	V4 = uuid.V4 // Version 4 (random)
	V5 = uuid.V5 // Version 5 (namespace name-based)
)

// UUID layout variants.
const (
	VariantNCS       = uuid.VariantNCS
	VariantRFC4122   = uuid.VariantRFC4122
	VariantMicrosoft = uuid.VariantMicrosoft
	VariantFuture    = uuid.VariantFuture
)

// UUID DCE domains.
const (
	DomainPerson = uuid.DomainPerson
	DomainGroup  = uuid.DomainGroup
	DomainOrg    = uuid.DomainOrg
)

// Timestamp is the count of 100-nanosecond intervals since 00:00:00.00,
// 15 October 1582 within a V1 UUID. This type has no meaning for V2-V5
// UUIDs since they don't have an embedded timestamp.
type Timestamp = uuid.Timestamp

// TimestampFromV1 returns the Timestamp embedded within a V1 UUID.
// Returns an error if the UUID is any version other than 1.
func TimestampFromV1(u UUID) (Timestamp, error) {
	return uuid.TimestampFromV1(u)
}

// Nil is the nil UUID, as specified in RFC-4122, that has all 128 bits set to
// zero.
var Nil = uuid.Nil

// Predefined namespace UUIDs.
var (
	NamespaceDNS  = uuid.NamespaceDNS
	NamespaceURL  = uuid.NamespaceURL
	NamespaceOID  = uuid.NamespaceOID
	NamespaceX500 = uuid.NamespaceX500
)

// Must is a helper that wraps a call to a function returning (UUID, error)
// and panics if the error is non-nil. It is intended for use in variable
// initializations such as
//  var packageUUID = v3.Must(v3.FromString("123e4567-e89b-12d3-a456-426655440000"))
func Must(u UUID, err error) UUID {
	return uuid.Must(u, err)
}

// NullUUID can be used with the standard sql package to represent a
// UUID value that can be NULL in the database.
type NullUUID = uuid.NullUUID

// HWAddrFunc is the function type used to provide hardware (MAC) addresses.
type HWAddrFunc = uuid.HWAddrFunc

// Generator provides an interface for generating UUIDs.
type Generator = uuid.Generator

// DefaultGenerator is the default UUID Generator used by this package.
var DefaultGenerator = uuid.DefaultGenerator

// NewV1 returns a UUID based on the current timestamp and MAC address.
func NewV1() (UUID, error) {
	return uuid.NewV1()
}

// NewV2 returns a DCE Security UUID based on the POSIX UID/GID.
func NewV2(domain byte) (UUID, error) {
	return uuid.NewV2(domain)
}

// NewV3 returns a UUID based on the MD5 hash of the namespace UUID and name.
func NewV3(ns UUID, name string) UUID {
	return uuid.NewV3(ns, name)
}

// NewV4 returns a randomly generated UUID.
func NewV4() (UUID, error) {
	return uuid.NewV4()
}

// NewV5 returns a UUID based on SHA-1 hash of the namespace UUID and name.
func NewV5(ns UUID, name string) UUID {
	return uuid.NewV5(ns, name)
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
type Gen = uuid.Gen

// NewGen returns a new instance of Gen with some default values set. Most
// people should use this.
func NewGen() *Gen {
	return uuid.NewGen()
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
	return uuid.NewGenWithHWAF(hwaf)
}

// FromBytes returns a UUID generated from the raw byte slice input.
// It will return an error if the slice isn't 16 bytes long.
func FromBytes(input []byte) (UUID, error) {
	return uuid.FromBytes(input)
}

// FromBytesOrNil returns a UUID generated from the raw byte slice input.
// Same behavior as FromBytes(), but returns v3.Nil instead of an error.
func FromBytesOrNil(input []byte) UUID {
	return uuid.FromBytesOrNil(input)
}

// FromString returns a UUID parsed from the input string.
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (UUID, error) {
	return uuid.FromString(input)
}

// FromStringOrNil returns a UUID parsed from the input string.
// Same behavior as FromString(), but returns v3.Nil instead of an error.
func FromStringOrNil(input string) UUID {
	return uuid.FromStringOrNil(input)
}
