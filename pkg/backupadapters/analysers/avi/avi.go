package avi

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"io"
	"path"
	"strings"
	"time"
)

type Parser struct {
	Debug bool
}

func (p *Parser) Supports(media backup.FoundMedia, mediaType backup.MediaType) bool {
	ext := strings.ToLower(path.Ext(media.MediaPath().Filename))
	return mediaType == backup.MediaTypeVideo && ext == ".avi"
}

func (p *Parser) ReadDetails(reader io.Reader, options backup.DetailsReaderOptions) (*backup.MediaDetails, error) {
	details := new(backup.MediaDetails)

	riffNode, err := DecodeRiff(reader)
	if err != nil {
		return nil, err
	}

	if p.Debug {
		printRiff(riffNode, 0)
	}

	browseRiffNodes(details, riffNode)

	return details, nil
}

func browseRiffNodes(details *backup.MediaDetails, node *RiffNode) {
	if node.Type == "IDIT" {
		value := string(node.Value)
		if len(value) > 24 {
			value = value[:24]
		}
		val, err := time.Parse("Mon Jan 02 15:04:05 2006", value) // ex: "Sun Sep 02 11:05:41 2007"
		if err != nil {
			log.Warnf("RIFF date parser failed for time %s: %s", value, err.Error())
		} else {
			details.DateTime = val
		}

	} else if node.Type == "strf" && len(node.Value) == 40 {
		details.Width = int(binary.LittleEndian.Uint32(node.Value[4:8]))
		details.Height = int(binary.LittleEndian.Uint32(node.Value[8:16]))

	} else if node.Type == "strh" && len(node.Value) == 56 {
		if binary.LittleEndian.Uint32(node.Value[4:8]) != 0 {
			details.VideoEncoding = strings.Trim(string(node.Value[4:8]), "\x00")
		}

		// note - likely to pick both VIDEO and AUDIO stream but duration is the same
		rate := float64(binary.LittleEndian.Uint32(node.Value[20:24])) / float64(binary.LittleEndian.Uint32(node.Value[24:28]))
		frameCount := float64(binary.LittleEndian.Uint32(node.Value[32:36]))

		details.Duration = int64(rate * frameCount * 1000)

	} else if node.Type == "ISFT" && len(node.Value) > 0 {
		details.Make = strings.Trim(string(node.Value), "\x00")
	}

	for _, child := range node.Children {
		browseRiffNodes(details, child)
	}
}

func printRiff(node *RiffNode, depth int) {
	fmt.Printf("%s- %s", strings.Repeat("  ", depth), node.Type)
	if len(node.Value) > 0 {
		fmt.Printf(" -> [%d bytes] %s\n", len(node.Value), hex.EncodeToString(node.Value))
	} else {
		fmt.Println()
	}

	for _, n := range node.Children {
		printRiff(n, depth+1)
	}
}
