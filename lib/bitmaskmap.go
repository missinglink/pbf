package lib

import (
	"encoding/gob"
	"io"
	"os"
)

// BitmaskMap - struct to hold common masks
type BitmaskMap struct {
	Nodes     *Bitmask
	Ways      *Bitmask
	Relations *Bitmask
}

// NewBitmaskMap - constructor
func NewBitmaskMap() *BitmaskMap {
	return &BitmaskMap{
		Nodes:     NewBitMask(),
		Ways:      NewBitMask(),
		Relations: NewBitMask(),
	}
}

// WriteTo - write to destination
func (m *BitmaskMap) WriteTo(sink io.Writer) (int64, error) {
	encoder := gob.NewEncoder(sink)
	err := encoder.Encode(m)
	return 0, err
}

// ReadFrom - read from destination
func (m *BitmaskMap) ReadFrom(tap io.Reader) (int64, error) {
	decoder := gob.NewDecoder(tap)
	err := decoder.Decode(m)
	return 0, err
}

// WriteToFile - write to disk
func (m *BitmaskMap) WriteToFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	m.WriteTo(file)
}

// ReadFromFile - read from disk
func (m *BitmaskMap) ReadFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	m.ReadFrom(file)
}
