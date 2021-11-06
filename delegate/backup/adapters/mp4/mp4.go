// Package mp4 parse a MP4 file to retrieve creation date, length, and other details found in the stream.
// References:
// - https://xhelmboyx.tripod.com/formats/mp4-layout.txt
// - https://www.programmersought.com/article/92132468003/
// - https://github.com/exiftool/exiftool/blob/master/lib/Image/ExifTool/QuickTime.pm
package mp4

import (
	"github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	Debug bool
}

func (p *Parser) Supports(media backupmodel.FoundMedia, mediaType backupmodel.MediaType) bool {
	ext := strings.ToUpper(path.Ext(media.MediaPath().Filename))
	return ext == ".MP4" || ext == ".MOV"
}

func (p *Parser) ReadDetails(reader io.Reader, options backupmodel.DetailsReaderOptions) (*backupmodel.MediaDetails, error) {
	details := new(backupmodel.MediaDetails)

	decoder := NewAtomDecoder(reader)
	for {
		atom, payload, err := decoder.Next()
		if err == io.EOF || options.Fast && !details.DateTime.IsZero() || !details.DateTime.IsZero() && details.VideoEncoding != "" && details.Width > 0 && details.Height > 0 && details.GPSLongitude != 0 && details.GPSLatitude != 0 {
			// end of file (EOF is expected to notify the end of the file) ; or all details already extracted
			break
		} else if err != nil {
			return nil, err
		}

		if p.Debug {
			fmt.Printf("'%s' atom [%d]:\n", atom.Path, atom.Size)
		}

		switch atom.Path {
		case "ftyp":
			buffer, err := payload.Next(4)
			if err != nil {
				return nil, err
			}

			details.VideoEncoding = strings.Trim(string(buffer), " ")
			if details.VideoEncoding == "mp42" {
				details.VideoEncoding = "MP4"
			}

		case "moov.mvhd":
			buffer, err := payload.Next(128)
			if err != nil {
				return nil, err
			}
			p.parseMVHD(details, buffer)

		case "moov.trak.tkhd":
			buffer, err := payload.Next(128)
			if err != nil {
				return nil, err
			}
			p.parseTKHD(details, buffer)

		case "moov.udta":
			buffer, err := payload.Next(1024)
			if err != nil {
				return nil, err
			}
			p.parseUDTA(details, buffer)

		default:
			if p.Debug {
				if !atom.IsParent {
					if payload.HasNext() {
						buffer, err := payload.Next(1024)
						if err != nil {
							return nil, err
						}

						fmt.Print(hex.Dump(buffer))
					}
					if payload.HasNext() {
						fmt.Println("...")
					}
				}
			}
		}
	}

	return details, nil
}

func (p *Parser) parseMVHD(details *backupmodel.MediaDetails, payload []byte) {
	version := payload[0]
	var timestampsFrom1904 uint64
	var timescale uint32
	var duration uint64

	if version >= 1 && len(payload) >= 32 {
		timestampsFrom1904 = binary.BigEndian.Uint64(payload[4:12])
		timescale = binary.BigEndian.Uint32(payload[20:24])
		duration = binary.BigEndian.Uint64(payload[24:32])
	} else if len(payload) >= 24 {
		timestampsFrom1904 = uint64(binary.BigEndian.Uint32(payload[4:8]))
		timescale = binary.BigEndian.Uint32(payload[12:16])
		duration = uint64(binary.BigEndian.Uint32(payload[16:24]))
	}

	if timestampsFrom1904 > 0 && timescale > 0 {
		details.DateTime = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(timestampsFrom1904) * time.Second)
		details.Duration = int64(1000 * duration / uint64(timescale))
	}
}

func (p *Parser) parseTKHD(details *backupmodel.MediaDetails, buffer []byte) {
	version := buffer[0]

	dateSize := 4
	if version >= 1 {
		dateSize = 8
	}
	// 1 (version) + 3 (unused after version) + 3*dateSize (creation, modification, and duration) + 4 (track id) + 8 (reserved) + 4 (reserved) + 2 (layer) + 2 (alternate group) + 2 (volume) + 2 (reserved) + 36 (video matrix)
	buffer = buffer[4+3*dateSize+24+36:]

	// long fixed point width, '2 bytes . 2 bytes'
	details.Width = int(binary.BigEndian.Uint16(buffer[:2]))
	details.Height = int(binary.BigEndian.Uint16(buffer[4:6]))
}

func (p *Parser) parseUDTA(details *backupmodel.MediaDetails, payload []byte) {
	if len(payload) >= 30 && binary.BigEndian.Uint32(payload[4:8]) == binary.BigEndian.Uint32([]byte("\xa9xyz")) {
		gps := string(payload[12:30])
		details.GPSLongitude, details.GPSLatitude = parseISO6709(gps)
	}
}

// parseISO6709 parse a ISO-6709 GPS coordinates
func parseISO6709(gps string) (float64, float64) {
	degMatcher := regexp.MustCompile("(?P<LON_DEG>[+-]\\d{1,2}(\\.\\d*)?)(?P<LAT_DEG>[+-]\\d{1,3}(\\.\\d*)?)")
	minutesMatcher := regexp.MustCompile("(?P<LON_DEG>[+-]\\d{1,2})(?P<LON_MIN>\\d{2}(\\.\\d*)?)(?P<LAT_DEG>[+-]\\d{1,3})(?P<LAT_MIN>\\d{2}(\\.\\d*)?)")
	secMatcher := regexp.MustCompile("(?P<LON_DEG>[+-]\\d{1,2})(?P<LON_MIN>\\d{2})(?P<LON_SEC>\\d{2}(\\.\\d*)?)(?P<LAT_DEG>[+-]\\d{1,3})(?P<LAT_MIN>\\d{2})(?P<LAT_SEC>\\d{2}(\\.\\d*)?)")

	submatch := degMatcher.FindStringSubmatch(gps)
	expNames := degMatcher.SubexpNames()
	if len(submatch) == 0 {
		submatch = minutesMatcher.FindStringSubmatch(gps)
		expNames = minutesMatcher.SubexpNames()
	}
	if len(submatch) == 0 {
		submatch = secMatcher.FindStringSubmatch(gps)
		expNames = secMatcher.SubexpNames()
	}
	if len(submatch) >= len(expNames) {
		var lonDeg, lonMin, lonSec float64
		var latDeg, latMin, latSec float64

		for i, key := range expNames {
			var err error

			switch key {
			case "LON_DEG":
				lonDeg, err = strconv.ParseFloat(submatch[i], 64)
			case "LON_MIN":
				lonMin, err = strconv.ParseFloat(submatch[i], 64)
			case "LON_SEC":
				lonSec, err = strconv.ParseFloat(submatch[i], 64)
			case "LAT_DEG":
				latDeg, err = strconv.ParseFloat(submatch[i], 64)
			case "LAT_MIN":
				latMin, err = strconv.ParseFloat(submatch[i], 64)
			case "LAT_SEC":
				latSec, err = strconv.ParseFloat(submatch[i], 64)
			}

			if err != nil {
				return 0, 0
			}
		}

		return lonDeg + lonMin/60 + lonSec/3600, latDeg + latMin/60 + latSec/3600
	}

	return 0, 0
}
