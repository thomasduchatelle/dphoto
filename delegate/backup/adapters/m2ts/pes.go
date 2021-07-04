package m2ts

import (
	"encoding/binary"
	"github.com/pkg/errors"
)

type PES []byte

// Valid checks if the payload starts with appropriate prefix (0x000001)
func (p PES) Valid() error {
	starts := binary.BigEndian.Uint32(p[0:4]) & 0xffffff00
	if len(p) < 6 || starts != 0x100 {
		return errors.Errorf("invalid PES payload: length must be more than 6 and it must start by 0x000001 (length=%d ; start=0x%x)", len(p), starts)
	}

	return nil
}

// PESLength is the number of bytes the PES header is. The Following bytes are the stream payload.
// Example:
//    00000000  00 00 01 e0 00 00 85 c0  0a 31 00 07 13 81 11 00  |.........1......|
//    00000010  05 bf 21 00 00 00 01 09  10 00 00 00 01 27 64 00  |..!..........'d.|
//                       ^ Video Stream starts here
func (p PES) PESLength() int {
	startCodeLength := 4
	pesPacketCounterLength := 2
	pesOptionalHeaderLength := 2
	optionalFieldsLength := int(p[startCodeLength+pesPacketCounterLength+pesOptionalHeaderLength]) // last byte of the PES optional header is the length of the optional fields

	return startCodeLength + pesOptionalHeaderLength + pesOptionalHeaderLength + optionalFieldsLength + 1
}
