// Package mp4 parse a MP4 file to retrieve creation date, length, and other details found in the stream.
// References:
// - https://xhelmboyx.tripod.com/formats/mp4-layout.txt
// - https://www.programmersought.com/article/92132468003/
package mp4

import (
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"path"
	"strings"
	"time"
)

type Parser struct {
	Debug bool
}

func (p *Parser) Supports(media backupmodel.FoundMedia, mediaType backupmodel.MediaType) bool {
	ext := strings.ToUpper(path.Ext(media.Filename()))
	return ext == "MP4"
}

func (p *Parser) ReadDetails(reader io.Reader, options backupmodel.DetailsReaderOptions) (*backupmodel.MediaDetails, error) {
	details := new(backupmodel.MediaDetails)

	decoder := NewAtomDecoder(reader)
	for {
		atom, payload, err := decoder.Next()
		if err == io.EOF || options.Fast && !details.DateTime.IsZero() || !details.DateTime.IsZero() && details.VideoEncoding != "" && details.Width > 0 && details.Height > 0 {
			// end of file (EOF is expected to notify the end of the file) ; or all details already extracted
			break
		} else if err != nil {
			return nil, err
		}

		switch atom.Path {
		case "ftyp":
			buffer, err := payload.Next(4)
			if err != nil {
				return nil, err
			}

			details.VideoEncoding = string(buffer)
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

		default:
			if p.Debug {
				if atom.IsParent {
					fmt.Printf("Starting a new '%s'\n", atom.Path)
				} else {
					fmt.Printf("'%s' atom:\n", atom.Path)
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

//func (e *Parser) ReadDetails2(reader io.Reader, options backupmodel.DetailsReaderOptions) (*backupmodel.MediaDetails, error) {
//	//buf := make([]byte, 2048)
//	//_, _ = io.ReadFull(reader, buf)
//	//fmt.Println(hex.Dump(buf))
//
//	firstChunk, err := readChunk(reader)
//	if err != nil {
//		return nil, err
//	}
//
//	if firstChunk.Code != "ftyp" {
//		return nil, errors.Errorf("not MP4 format, first chunk type expected to be 0x66747970 but is %x", firstChunk.Code)
//	}
//
//	fileMetadata, _, err := firstChunk.Payload(reader)
//	if err != nil {
//		return nil, err
//	}
//
//	details := new(backupmodel.MediaDetails)
//	details.VideoEncoding = string(fileMetadata[:4])
//	if details.VideoEncoding == "mp42" {
//		details.VideoEncoding = "MP4"
//	}
//
//	moov, err := readChunk(reader)
//	if err != nil {
//		return nil, err
//	}
//	for moov.Code != "moov" {
//		// search for 'moov' atom, skip any previous one
//		var more bool
//		for _, more, err = moov.Payload(reader); err == nil && more; _, more, err = moov.Payload(reader) {
//			// read whole payload...
//		}
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	moreChildren := moov.Size > 8
//	for moreChildren {
//		var child *Atom
//		child, moreChildren, err = moov.ReadChild(reader)
//		if err != nil {
//			return nil, err
//		}
//
//		switch child.Code {
//		case "mvhd":
//			payload, _, err := child.Payload(reader)
//			if err != nil {
//				return nil, err
//			}
//			e.parseMVHD(payload, details)
//
//			if options.Fast {
//				// In Fast mode, only datetime matter, return as soon it has been retrieved
//				moreChildren = false
//			}
//
//		case "trak":
//			err = e.parseTRACK(child, reader, details)
//			if err != nil {
//				return nil, err
//			}
//
//		default:
//			payload, more, err := child.Payload(reader)
//			if err != nil {
//				return nil, err
//			}
//
//			if e.Debug {
//				fmt.Printf("\nChunk %s:\n%s", child.Code, hex.Dump(payload))
//			}
//
//			if more {
//				// if chunk is too big, there is certainly no more metadata to be found
//				fmt.Println("...")
//				moreChildren = false
//			}
//
//		}
//	}
//
//	return details, nil
//}

//func (e *Parser) parseTRACK(trak *Atom, reader io.Reader, details *backupmodel.MediaDetails) (err error) {
//	fmt.Println("Parsing track...")
//	hasNext := trak.Size > 8
//	for hasNext {
//		var child *Atom
//		child, hasNext, err = trak.ReadChild(reader)
//		if err != nil {
//			return err
//		}
//		switch child.Code {
//		case "mdia":
//
//		default:
//			payload, _, err := child.Payload(reader)
//			if err != nil {
//				return err
//			}
//
//			fmt.Printf("\nChunk %s:\n%s", child.Code, hex.Dump(payload))
//		}
//
//	}
//
//	return nil
//}

//func (c *Atom) Payload(reader io.Reader) ([]byte, bool, error) {
//	bufferSize := c.Size - c.read
//	if bufferSize > 1024 {
//		bufferSize = 1024
//	}
//
//	buffer := make([]byte, bufferSize)
//	_, err := io.ReadFull(reader, buffer)
//
//	c.read += bufferSize
//	return buffer, c.read < c.Size, err
//}
//
//func (c *Atom) ReadChild(reader io.Reader) (*Atom, bool, error) {
//	if c.read >= c.Size {
//		return nil, false, errors.Errorf("There is no more child in atom %+v", c)
//	}
//
//	child, err := readChunk(reader)
//	if err != nil {
//		return nil, false, err
//	}
//
//	c.read += child.Size
//	return child, c.read < c.Size, nil
//}
