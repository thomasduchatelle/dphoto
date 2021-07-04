package m2ts

import (
	"encoding/binary"
	"github.com/pkg/errors"
)

// PSI is a wrapper to decode a PSI table
type PSI []byte

const (
	mask10bits = 0x03ff // mask10bits is a mask to apply to any integer to only keep last 10 bits
	crcLength  = 4      // crcLength is the number of bytes of CRC check sum in PSI table
)

func (p PSI) Valid() error {
	if len(p) < 4 || p[2]&0xbc != 0xb0 {
		return errors.Errorf("invalid PSI content: length must be more than 4 and bits 19-20 must be 1 (len=%d ; 3rd byte = %x)", len(p), p[2]&0xbc)
	}

	return nil
}

func (p PSI) sectionLength() uint16 {
	return binary.BigEndian.Uint16(p[2:4]) & 0x03ff
}

func (p PSI) tableId() byte {
	return p[1]
}

// AssociationTable reads PAT table content to extract association table
func (p PSI) AssociationTable() map[uint16]uint16 {
	programs := make(map[uint16]uint16)

	if p.tableId() == 0x0 {
		firstPAT := 9
		maxIndex := int(4 + p.sectionLength() - crcLength)
		for i := firstPAT; i < maxIndex && i+4 <= len(p); i += 4 {
			program := binary.BigEndian.Uint16(p[i : i+2])
			pid := binary.BigEndian.Uint16(p[i+2:i+4]) & PIDMask

			programs[pid] = program
		}
	}

	return programs
}

// UpdateProgramMap reads PMT table to update program map
func (p PSI) UpdateProgramMap(programs map[uint16]ProgramNumber) map[uint16]ProgramNumber {
	pmtStart := 9
	maxIndex := int(4 + p.sectionLength() - crcLength)

	streamIdx := pmtStart + 4 + int(binary.BigEndian.Uint16(p[pmtStart+2:pmtStart+4])&mask10bits)

	for streamIdx < maxIndex {
		streamType := p[streamIdx]
		pid := binary.BigEndian.Uint16(p[streamIdx+1:streamIdx+3]) & PIDMask
		programs[pid] = ProgramNumber(streamType)

		entryLength := 5 + int(binary.BigEndian.Uint16(p[streamIdx+3:streamIdx+5])&mask10bits)
		streamIdx += entryLength
	}

	return programs
}
