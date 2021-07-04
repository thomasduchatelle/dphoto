package m2ts

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/pkg/errors"
)

const (
	extraHeader = 4
	syncByte    = 0x47
)

type Packet []byte

// NewPacket creates a Packet structure to parse M2TS files
func NewPacket() []byte {
	return make([]byte, 192)
}

// Valid confirms it's a M2TS packet by looking at the Sync Byte
func (p Packet) Valid() error {
	if p[extraHeader] != syncByte {
		return errors.Errorf("Invalid M2TS (.MTS) packet: sync byte (0x47) is missing [%s]", hex.EncodeToString(p))
	}
	if len(p) < syncByte+4 {
		return errors.Errorf("Invalid M2TS (.MTS) packat: it must be at least 8 bytes [%s]", hex.EncodeToString(p))
	}
	return nil
}

// PID returns packet PID ; it might panic if packet is not Valid()
func (p Packet) PID() uint16 {
	return binary.BigEndian.Uint16(p[extraHeader+1:extraHeader+3]) & 0x1fff
}

// PCR retrieve PCR value in adaptation field if present, or returns (0, false)
func (p Packet) PCR() (uint64, bool) {
	if p.hasAdaptationField() {
		adaptationFieldLength := p[extraHeader+4]
		if adaptationFieldLength >= 7 && (p[extraHeader+5]&0x10) != 0 {
			pcrBuf := p[10:16]
			base := uint64(binary.BigEndian.Uint32(pcrBuf[:4]))<<1 + uint64((pcrBuf[4]&0x80)/0x80)
			ext := uint64(pcrBuf[5]) + uint64(pcrBuf[4]&0x1)*0x100

			return base*300 + ext, true
		}
	}

	return 0, false
}

func (p Packet) hasAdaptationField() bool {
	return p.adaptationFieldControl()&0x20 != 0
}

func (p Packet) adaptationFieldControl() byte {
	return p[extraHeader+3] & 0x30
}

// Payload retrieves the payload from the packet or returns nil
func (p Packet) Payload() []byte {
	if p.adaptationFieldControl()&0x10 == 0 {
		return nil
	}

	start := extraHeader + 4
	if p.hasAdaptationField() {
		start += int(p[extraHeader+4]) + 1
	}

	return p[start:]
}

// PesPsiDvbPayload returns TRUE when the Payload is PES, PSI, or DVB
func (p Packet) PesPsiDvbPayload() bool {
	return p[extraHeader+1]&0x40 != 0
}
