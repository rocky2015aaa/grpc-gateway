package util

import (
	"encoding/binary"
	"errors"
	"os"

	"github.com/alfg/mp4/atom"
)

func Mp4Duration(path string) (float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	f := &atom.File{File: file, Size: info.Size()}

	for offset := int64(0); offset < f.Size; {
		boxSize, boxType := f.ReadBoxAt(offset)
		if boxType == "moov" {
			offset += atom.BoxHeaderSize

			for offset < (f.Size - atom.BoxHeaderSize) {
				boxSize, boxType = f.ReadBoxAt(offset)
				if boxType == "mvhd" {
					mvhd := &atom.Box{
						Name:  string(boxType),
						Size:  int64(boxSize),
						File:  f,
						Start: offset,
					}
					data := mvhd.ReadBoxData()

					ts := binary.BigEndian.Uint32(data[12:16])
					duration := binary.BigEndian.Uint32(data[16:20])
					return float64(duration) / float64(ts), nil
				}

				offset += int64(boxSize)
			}
		}

		offset += int64(boxSize)
	}
	return 0, errors.New("failed to get mp4 duration")
}
