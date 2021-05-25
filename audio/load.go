package audio

import (
	"encoding/binary"
	"io"
	"os"
)

// LoadAudio attempts to load an encoded sound file from disk.
func LoadAudio(f string, b *[][]byte) error {
	file, err := os.Open(f)
	if err != nil {
		return err
	}

	var audioLength int16
	for {
		err = binary.Read(file, binary.LittleEndian, &audioLength)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}

		InBuf := make([]byte, audioLength)
		err = binary.Read(file, binary.LittleEndian, &InBuf)
		if err != nil {
			return err
		}

		// Append encoded pcm data to the buffer.
		*b = append(*b, InBuf)
	}
}
