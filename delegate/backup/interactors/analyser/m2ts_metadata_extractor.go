package analyser

import (
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"fmt"
	"io"
)

func m2tsAdapter(found backupmodel.FoundMedia) (*backupmodel.MediaDetails, error) {
	reader, err := found.ReadMedia()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 192)
	for read, err := reader.Read(buf); read > 0; {
		if err != nil && err != io.EOF {
			return nil, err
		}

		if read >= 9 {
			syncByte := buf[8]
			ok := syncByte == 0x47
			fmt.Printf("Sync byte: %b -> %t", syncByte, ok)
		}
	}

	return nil, err
}
