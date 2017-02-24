package gosmparse

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
)

// GroupInfo - store info about each group
type GroupInfo struct {
	Type  string
	Count int
	High  int64
	Low   int64
}

// BlobInfo - store info about each block
type BlobInfo struct {
	Groups []*GroupInfo
	Start  uint64
	Size   uint64
}

// BlobIndex - an index of all blocks in the file
type BlobIndex struct {
	Blobs []*BlobInfo
}

// BlobOffsets - find the start offset of blob(s) containing desired element
func (i *BlobIndex) BlobOffsets(memtype string, id int64) ([]int64, error) {
	var offsets []int64
	for _, info := range i.Blobs {
		for _, group := range info.Groups {
			if group.Type == memtype {
				if id <= group.High && id >= group.Low {
					offsets = append(offsets, int64(info.Start))
				}
			}
		}
	}

	if len(offsets) > 0 {
		return offsets, nil
	}

	return offsets, errors.New("not found")
}

// WriteTo - write to destination
func (i *BlobIndex) WriteTo(sink io.Writer) (int64, error) {
	encoder := gob.NewEncoder(sink)
	err := encoder.Encode(i)
	return 0, err
}

// ReadFrom - read from destination
func (i *BlobIndex) ReadFrom(tap io.Reader) (int64, error) {
	decoder := gob.NewDecoder(tap)
	err := decoder.Decode(i)
	return 0, err
}

// WriteToFile - write to disk
func (i *BlobIndex) WriteToFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	i.WriteTo(file)
}

// ReadFromFile - read from disk
func (i *BlobIndex) ReadFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	i.ReadFrom(file)
}
