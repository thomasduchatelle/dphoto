package m2ts

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

const (
	mdpmEntrySize    = 5
	MDPMDate         = 0x18
	MDPMTime         = 0x19
	MDPMGPSLatitude  = 0xb2
	MDPMGPSLongitude = 0xb6
	MDPMMake         = 0xe0
	MDPMModel        = 0xe4
)

var (
	mdpmToken = []byte{0x17, 0xee, 0x8c, 0x60, 0xf8, 0x4d, 0x11, 0xd9, 0x8c, 0xd6, 0x08, 0, 0x20, 0x0c, 0x9a, 0x66, 'M', 'D', 'P', 'M'}
	makeMap   = map[uint16]string{
		0x0103: "Panasonic",
		0x0108: "Sony",
		0x1011: "Canon",
		0x1104: "JVC",
	}
)

// UpdateDetailsFromMDPM find MDPM UUID followed by 'MDPM' marker, and reads what's following as key-value map
// inspired by https://metacpan.org/release/EXIFTOOL/Image-ExifTool-8.90/source/lib/Image/ExifTool/H264.pm
func UpdateDetailsFromMDPM(payload []byte, details *backup.MediaDetails) map[byte][]byte {
	mdpm := readMDPM(payload)

	details.DateTime, _ = extractDateFromMDPM(mdpm)

	if makeValue, ok := mdpm[MDPMMake]; ok {
		details.Make, _ = makeMap[binary.BigEndian.Uint16(makeValue[:2])]
	}

	details.Model = extractModelFromMDPM(mdpm)

	return mdpm
}

func extractModelFromMDPM(mdpm map[byte][]byte) string {
	// Model is only working for SONY
	var modelBuffer []byte
	ended := false
	for key := byte(0xe4); !ended && key < 0xe7; key++ {
		val, ok := mdpm[key]
		ended = ended || !ok
		for i := 0; !ended && i < len(val); i++ {
			if val[i] == 0 {
				ended = true
			} else if val[i] <= unicode.MaxASCII {
				modelBuffer = append(modelBuffer, val[i])
			}
		}
	}

	return string(modelBuffer)
}

func extractDateFromMDPM(mdpm map[byte][]byte) (time.Time, error) {

	var dateTime []byte
	dateBuf, hasDate := mdpm[MDPMDate]
	timeBuf, hasTime := mdpm[MDPMTime]
	dateTime = append(dateTime, dateBuf...)
	dateTime = append(dateTime, timeBuf...)

	if hasDate && hasTime && len(dateTime) == 8 {
		// Time zone first byte
		//  0x80 - unused
		//  0x40 - DST flag
		//  0x20 - TimeZoneSign
		//  0x1e - TimeZoneValue
		//  0x01 - half-hour flag
		sign := "+"
		hours := int(0xf & (dateTime[0] >> 1))
		minutes := "00"
		if dateTime[0]&0x02 == 0 {
			sign = "-"
		}
		if dateTime[0]&0x1 != 0 {
			minutes = "30"
		}

		return time.Parse("20060102150405-0700", fmt.Sprintf("%x%s%02d%s", dateTime[1:], sign, hours, minutes))
	}

	return time.Time{}, errors.Errorf("No date-time found in MDPM %s", printMDPM(mdpm))
}

func readMDPM(payload []byte) map[byte][]byte {
	mdpm := make(map[byte][]byte)

	start := 0
	for start < len(payload) && !bytes.HasPrefix(payload[start:], mdpmToken) {
		start++
	}

	start += len(mdpmToken) + 1
	if start >= len(payload) {
		return mdpm
	}

	size := int(payload[start] - 1)
	for i := 0; i < size && start+i*mdpmEntrySize+5 < len(payload); i++ {
		// each entry is 5 bytes: 1 for the key, 4 for the value
		entryIndex := start + i*mdpmEntrySize
		key := payload[entryIndex]
		value := payload[entryIndex+1 : entryIndex+5]
		mdpm[key] = value
	}

	return mdpm
}

func printMDPM(mdpm map[byte][]byte) string {
	mdpmDump := []string{
		"MDPM (H264 Netadata)",
	}

	keys := make([]byte, 0, len(mdpm))
	for k := range mdpm {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		mdpmDump = append(mdpmDump, fmt.Sprintf("\t- %02x = %s", k, regexp.MustCompile(" +").ReplaceAllString(strings.Trim(hex.Dump(mdpm[k]), "\n")[10:], " ")))
	}

	return strings.Join(mdpmDump, "\n")
}
