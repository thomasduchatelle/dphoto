// Package m2ts is parsing M2TS files followings specs:
// https://en.wikipedia.org/wiki/MPEG_transport_stream
// https://en.wikipedia.org/wiki/Packetized_elementary_stream
// https://en.wikipedia.org/wiki/Program-specific_information
package m2ts

import (
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"path"
	"strings"
)

const (
	PIDMask                          = 0x1fff // PIDMask can be used to get PID from uint16: last 13 bits
	ProgramNumberH264  ProgramNumber = 0x1b   // ProgramNumberH264 is a H264 video stream
	minSizeH264Payload               = 256
)

var (
	packetGroups = []int{4, 1, 3}
	psiGroups    = []int{4, 5, 4, 4, 4, 4, 4, 4, 4, 4, 4}
)

// ProgramNumber is the code as defined in https://en.wikipedia.org/wiki/Program-specific_information#Program_and_Elementary_Stream_Descriptor_Tags
type ProgramNumber uint8

type Parser struct {
	Debug bool // Debug can be set to true to print in the console payload dumps and read metadata.
}

func (p *Parser) Supports(media backupmodel.FoundMedia, mediaType backupmodel.MediaType) bool {
	ext := strings.ToUpper(path.Ext(media.Filename()))
	return ext == "MTS" || ext == "M2TS"
}

// ReadDetails unmux M2TS (MTS) file, with h264 support, to collect the Make, Model, and DateTime of the video flux.
// Example:
// 00 00 00 00 | 47 40 00 10 | 00 00 b0 11 00 00 c1 00 00 00 00 e0 1f 00 01 e1 00 23 5a ab 82 ffffff ffffffff ffffffff ffffffff ffffffff ffffffff ffffffff ffffffff ffffffff
//               M2TS HEADER-- PAT HEADER- SECTION------- #1--------- #2--------- CRC--------
// 00 00 06 9c | 47 | 41 00 10 | 00 02 b0 3e | 00 01 c1 00 00 | f0 01 f0 0c | 05 04 48 44 | 4d 56 88 04 | 0f ff fc fc | 1b f0 11 f0 0a | 05 08 48 44 4d 56 ff 1b 43 3f | 81 f1 00 f0 0c | 05 04 41 43 2d 33 81 04 04 30 04 00 90 f2 00 f0 00 | 0c d3 f4 dc
//               M2TS HEADER---  PAT HEADER--  SECTION-------   PMT-------- | Program Descriptors -------------------   #1 1011=>1b-----------------------------------                                                                         CRC--------
func (p *Parser) ReadDetails(reader io.Reader, options backupmodel.DetailsReaderOptions) (*backupmodel.MediaDetails, error) {
	details := new(backupmodel.MediaDetails)

	var packet Packet = make([]byte, 192)
	count := 0

	var minPCR, maxPCR uint64
	programs := make(map[uint16]ProgramNumber)
	payloads := make(map[uint16][]byte)
	payloadsFull := false

	for !options.Fast || !payloadsFull {
		count++

		full, err := io.ReadFull(reader, packet)
		if full < 192 || err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		err = packet.Valid()
		if err != nil {
			return nil, err
		}

		pid := packet.PID()
		if pcr, ok := packet.PCR(); ok {
			if minPCR == 0 {
				minPCR = pcr
			}
			maxPCR = pcr
		}

		payload := packet.Payload()
		if packet.PesPsiDvbPayload() {
			pes := PES(payload)
			pesErr := pes.Valid()

			psi := PSI(payload)
			psiErr := psi.Valid()

			switch {
			case pesErr == nil:
				// New streams starts by the PES header, remove it from the payload
				payload = payload[pes.PESLength():]

			case psiErr == nil:
				programs = psi.UpdateProgramMap(programs)

			default:
				return nil, errors.Wrapf(psiErr, "couldn't parse packet %d \n[%s]\n", count, bytesToString(packet, packetGroups, psiGroups))
			}
		}

		// H264 payload is where Make, Model, and DateTime are found. Only keeping 256 bytes.
		if streamType, ok := programs[pid]; ok && streamType == ProgramNumberH264 {
			buffer, _ := payloads[pid]
			if len(buffer) < minSizeH264Payload {
				payloads[pid] = append(buffer, payload...)
			} else {
				payloadsFull = true
			}
		}
	}

	// debug - print program map and payloads
	if p.Debug {
		fmt.Printf("%d programs found:\n", len(programs))
		for pid, programNumber := range programs {
			fmt.Printf("\t- 0x%04x -> 0x%02x\n", pid, programNumber)
		}

		for pid, payload := range payloads {
			fmt.Printf("PID %04x\n%s", pid, hex.Dump(payload))
		}
	}

	// PCR tick at 27MHz
	details.Duration = int64((maxPCR - minPCR) / 27000)

	// debug - H264 MDPM
	for pid, programNumber := range programs {
		if payload, ok := payloads[pid]; ok && programNumber == ProgramNumberH264 {
			details.VideoEncoding = "H264"

			mdpm := UpdateDetailsFromMDPM(payload, details)
			if p.Debug {
				fmt.Println(printMDPM(mdpm))
			}
		}
	}

	return details, nil
}

func bytesToString(buf []byte, separatorSlices ...[]int) string {
	separators := make([]int, 0)
	for _, sl := range separatorSlices {
		separators = append(separators, sl...)
	}

	blocks := make([]string, 0, len(buf)+len(separators))
	separatorIndex := 0
	separatorCount := 0

	for index, b := range buf {
		if separatorIndex < len(separators) && separators[separatorIndex]+separatorCount == index {
			blocks = append(blocks, "|")
			separatorCount += separators[separatorIndex]
			separatorIndex++
		}
		blocks = append(blocks, hex.EncodeToString([]byte{b}))
	}

	return strings.Join(blocks, " ")
}
