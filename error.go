package uuid

type Error string

const (
	ErrInvalidFormat           = Error("uuid: invalid UUID format")
	ErrIncorrectFormatInString = Error("uuid: incorrect UUID format in string")
	ErrIncorrectLength         = Error("uuid: incorrect UUID length")
	ErrIncorrectByteLength     = Error("uuid: UUID must be exactly 16 bytes long")
	ErrNoHwAddressFound        = Error("uuid: no HW address found")
	ErrTypeConvertError        = Error("uuid: cannot convert")
	ErrInvalidVersion          = Error("uuid:")
)

func (e Error) Error() string {
	return string(e)
}

func (e Error) Is(target error) bool {
	_, ok := target.(*Error)
	return ok
}
