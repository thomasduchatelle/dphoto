package mp4

import (
	"encoding/binary"
	"io"
	"strings"
)

var (
	parents = map[string]interface{}{
		"moov": nil,
		"trak": nil,
		"mdia": nil,
	}
)

// Atom is a chunk in a MP4 files, structured as a tree
type Atom struct {
	Path     string // Path is the dot separated code of the parents and of this atom (ex: moov.trak.hvdt)
	Code     string // Code is the 4 chars code used in MP4 specifications
	IsParent bool   // IsParent is TRUE when the Atom type is reconised to have children. The payload should not be read in this case.
	read     uint64
	Size     uint64 // Size is the total size of the Atom, including children (in bytes)
}

// AtomDecoder reads an MP4 file to extract Atom
type AtomDecoder struct {
	reader io.Reader
	stack  []*Atom
}

// PayloadIterator is used to get content of the Atom as a []byte
type PayloadIterator struct {
	reader io.Reader
	atom   *Atom
}

// NewAtomDecoder will start parsing an MP4 files
func NewAtomDecoder(reader io.Reader) *AtomDecoder {
	return &AtomDecoder{
		reader: reader,
	}
}

// Next reads and returns the next Atom (or children Atom). It will consume any payload not yet read.
// It will return (nil, nil, io.EOF) when the file has been read fully.
func (d *AtomDecoder) Next() (*Atom, *PayloadIterator, error) {
	for i := len(d.stack) - 1; i >= 0; i-- {
		atom := d.stack[i]

		// burning the payload if no children
		if !atom.IsParent {
			it := PayloadIterator{
				reader: d.reader,
				atom:   atom,
			}
			for it.HasNext() {
				_, err := it.Next(1024)
				if err != nil {
					return nil, nil, err
				}
			}
		}

		// removing atom if fully consumed
		if atom.read >= atom.Size {
			d.stack = d.stack[:i]
			if i > 0 {
				d.stack[i-1].read += atom.Size
			}
		}
	}

	atom, err := d.readChunk(d.reader)
	if err != nil {
		return nil, nil, err
	}

	d.stack = append(d.stack, atom)

	return atom, &PayloadIterator{
		reader: d.reader,
		atom:   atom,
	}, err
}

func (d *AtomDecoder) readChunk(reader io.Reader) (*Atom, error) {
	header := make([]byte, 8)
	_, err := io.ReadFull(reader, header)
	if err != nil {
		return nil, err
	}

	code := string(header[4:])
	_, isParent := parents[code]

	atomPath := make([]string, len(d.stack)+1)
	for i, c := range d.stack {
		atomPath[i] = c.Code
	}
	atomPath[len(d.stack)] = code

	size := uint64(binary.BigEndian.Uint32(header[:4]))
	read := uint64(8)

	if size == 1 {
		sizeBuffer := make([]byte, 8)
		_, err = io.ReadFull(reader, sizeBuffer)
		if err != nil {
			return nil, err
		}

		size = binary.BigEndian.Uint64(sizeBuffer)
		read = uint64(16)
	}

	return &Atom{
		Path:     strings.Join(atomPath, "."),
		Code:     code,
		IsParent: isParent,
		read:     read,
		Size:     size,
	}, err
}

// HasNext returns TRUE if some payload remains within the Atom.
func (p *PayloadIterator) HasNext() bool {
	return p.atom.read < p.atom.Size
}

// Next will read the Atom payload if any.
// Beware: it will read children as payload which might break AtomDecoder if used on an improper Atom.
func (p *PayloadIterator) Next(size int) ([]byte, error) {
	if size > int(p.atom.Size-p.atom.read) {
		size = int(p.atom.Size - p.atom.read)
	}

	buf := make([]byte, size)
	_, err := io.ReadFull(p.reader, buf)
	p.atom.read += uint64(len(buf))

	return buf, err
}
