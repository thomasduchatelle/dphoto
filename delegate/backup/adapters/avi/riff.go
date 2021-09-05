package avi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

type RiffNode struct {
	Type     string
	Value    []byte
	Children []*RiffNode
}

// IsRiff is looking at the first bytes to determine if it's RIFF structure (and not ATOM).
func IsRiff(reader io.Reader) (io.Reader, bool, error) {
	flag := make([]byte, 4)
	_, err := io.ReadFull(reader, flag)
	return io.MultiReader(bytes.NewReader(flag), reader), string(flag) == "RIFF", err
}

// DecodeRiff read the RIFF header, if any, and returns non-nil RiffNode and reader positioned to the MP4 content ; or nil RiffNode and unchanged reader
func DecodeRiff(reader io.Reader) (*RiffNode, error) {
	flag := make([]byte, 12)
	_, err := io.ReadFull(reader, flag)
	if err != nil {
		return nil, err
	}

	if string(flag[:4]) != "RIFF" {
		return nil, errors.Errorf("RIFF structure must start by 'RIFF', not '%s'", string(flag[:4]))
	}

	header := make([]byte, 12)
	_, err = io.ReadFull(reader, header)
	if err != nil {
		return nil, err
	}
	if string(header[:4]) != "LIST" {
		return nil, errors.Errorf("First chunk of RIFF is expected to be LIST but is %s", string(header[:4]))
	}

	length := binary.LittleEndian.Uint32(header[4:8]) - 4
	content := make([]byte, length)

	_, err = io.ReadFull(reader, content)
	if err != nil {
		return nil, err
	}

	nodes, err := parseRiffNode(content, 0)
	return &RiffNode{
		Type:     fmt.Sprintf("%s_%s", string(flag[8:]), string(header[8:])),
		Children: nodes,
	}, err
}

func parseRiffNode(content []byte, absIndex int) ([]*RiffNode, error) {
	if len(content) < 8 {
		return nil, errors.Errorf("RIFF entry should be at least 8 bytes [index: %d].", absIndex)
	}

	var nodes []*RiffNode
	i := 0
	for i < len(content) {
		name := string(content[i : i+4])
		size := int(binary.LittleEndian.Uint32(content[i+4 : i+8]))

		if len(content) < i+size {
			return nil, errors.Errorf("RIFF Chunk size is bigger than available data [index %d]: content's size = %d ; chunk size = %d", absIndex+i, len(content), i+size)
		}

		if name == "LIST" {
			listName := string(content[i+8 : i+12])
			listContent := content[i+12 : i+8+size]
			children, err := parseRiffNode(listContent, absIndex+i)
			if err != nil {
				return nil, err
			}

			nodes = append(nodes, &RiffNode{
				Type:     listName,
				Children: children,
			})

		} else {
			value := content[i+8 : i+8+size]
			nodes = append(nodes, &RiffNode{
				Type:  name,
				Value: value,
			})
		}

		i += 8 + int(size)
	}

	return nodes, nil
}
