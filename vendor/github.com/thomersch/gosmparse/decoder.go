package gosmparse

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/thomersch/gosmparse/OSMPBF"

	"github.com/golang/protobuf/proto"
)

// A Decoder reads and decodes OSM data from an input stream.
type Decoder struct {
	// QueueSize allows to tune the memory usage vs. parse speed.
	// A larger QueueSize will consume more memory, but may speed up the parsing process.
	QueueSize int
	r         *os.File
	o         OSMReader
	Mutex     *sync.Mutex
	BytesRead uint64
	Index     *BlobIndex
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r *os.File) *Decoder {
	return &Decoder{
		r:         r,
		QueueSize: 200,
	}
}

// ParseBlob - parse a single blob
func (d *Decoder) ParseBlob(o OSMReader, offset int64) error {

	// ---
	d.Index = &BlobIndex{}
	d.Mutex = &sync.Mutex{}

	d.o = o
	d.r.Seek(offset, 0)

	_, blob, err := d.block()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	err = d.readElements(blob)
	if err != nil {
		return err
	}

	return nil
}

// Parse starts the parsing process that will stream data into the given OSMReader.
func (d *Decoder) Parse(o OSMReader) error {

	d.Index = &BlobIndex{}
	d.Mutex = &sync.Mutex{}

	d.o = o
	header, _, err := d.block()
	if err != nil {
		return err
	}
	// TODO: parser checks
	if header.GetType() != "OSMHeader" {
		return fmt.Errorf("Invalid header of first data block. Wanted: OSMHeader, have: %s", header.GetType())
	}

	errChan := make(chan error)
	// feeder
	blobs := make(chan *OSMPBF.Blob, d.QueueSize)
	go func() {
		defer close(blobs)
		for {
			_, blob, err := d.block()
			if err != nil {
				if err == io.EOF {
					return
				}
				errChan <- err
				return
			}
			blobs <- blob
		}
	}()

	consumerCount := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	for i := 0; i < consumerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for blob := range blobs {
				err := d.readElements(blob)
				if err != nil {
					errChan <- err
					return
				}
			}
		}()
	}

	finished := make(chan bool)
	go func() {
		wg.Wait()
		finished <- true
	}()
	select {
	case err = <-errChan:
		return err
	case <-finished:
		return nil
	}
}

func (d *Decoder) block() (*OSMPBF.BlobHeader, *OSMPBF.Blob, error) {

	// store info
	var startBytes = d.BytesRead

	var byteCount int
	var err error

	// BlobHeaderLength
	headerSizeBuf := make([]byte, 4)

	// read BlobHeaderLength
	byteCount, err = io.ReadFull(d.r, headerSizeBuf)

	// keep track of bytes read so far
	atomic.AddUint64(&d.BytesRead, uint64(byteCount))

	// error checking
	if err != nil {
		return nil, nil, err
	}
	headerSize := binary.BigEndian.Uint32(headerSizeBuf)

	// BlobHeader
	headerBuf := make([]byte, headerSize)

	// read BlobHeader
	byteCount, err = io.ReadFull(d.r, headerBuf)

	// keep track of bytes read so far
	atomic.AddUint64(&d.BytesRead, uint64(byteCount))

	if err != nil {
		return nil, nil, err
	}

	blobHeader := new(OSMPBF.BlobHeader)
	if err = proto.Unmarshal(headerBuf, blobHeader); err != nil {
		return nil, nil, err
	}

	// Blob
	blobBuf := make([]byte, blobHeader.GetDatasize())
	byteCount, err = io.ReadFull(d.r, blobBuf)

	// keep track of bytes read so far
	atomic.AddUint64(&d.BytesRead, uint64(byteCount))

	if err != nil {
		return nil, nil, err
	}
	blob := new(OSMPBF.Blob)
	if err := proto.Unmarshal(blobBuf, blob); err != nil {
		return nil, nil, err
	}

	// store info
	d.Mutex.Lock()
	d.Index.Blobs = append(d.Index.Blobs, &BlobInfo{
		Start: startBytes,
		Size:  uint64(byteCount),
	})

	// hack to store the blob index key
	var key = make([]byte, 8)
	binary.LittleEndian.PutUint64(key, uint64(len(d.Index.Blobs)-1))
	blob.XXX_unrecognized = key

	d.Mutex.Unlock()

	return blobHeader, blob, nil
}

func (d *Decoder) readElements(blob *OSMPBF.Blob) error {

	pb, err := d.blobData(blob)
	if err != nil {
		return err
	}

	for _, pg := range pb.Primitivegroup {

		var info = &GroupInfo{}

		switch {
		case pg.Dense != nil:

			info.Type = "node"
			info.Count = len(pg.Dense.Id)

			// find high and low id
			var id int64
			for index := range pg.Dense.Id {
				id = pg.Dense.Id[index] + id
				if 0 == info.High || id > info.High {
					info.High = id
				}
				if 0 == info.Low || id < info.Low {
					info.Low = id
				}
			}

			if err := denseNode(d.o, pb, pg.Dense); err != nil {
				return err
			}
		case len(pg.Ways) != 0:

			info.Type = "way"
			info.Count = len(pg.Ways)

			// find high and low id
			var id int64
			for _, way := range pg.Ways {
				id = way.GetId()
				if 0 == info.High || id > info.High {
					info.High = id
				}
				if 0 == info.Low || id < info.Low {
					info.Low = id
				}
			}

			if err := way(d.o, pb, pg.Ways); err != nil {
				return err
			}
		case len(pg.Relations) != 0:

			info.Type = "relation"
			info.Count = len(pg.Relations)

			// find high and low id
			var id int64
			for _, way := range pg.Relations {
				id = way.GetId()
				if 0 == info.High || id > info.High {
					info.High = id
				}
				if 0 == info.Low || id < info.Low {
					info.Low = id
				}
			}

			if err := relation(d.o, pb, pg.Relations); err != nil {
				return err
			}
		case len(pg.Nodes) != 0:
			return fmt.Errorf("Nodes are not supported")
		default:
			return fmt.Errorf("no supported data in primitive group")
		}

		d.Mutex.Lock()

		// hack to retrieve the blob index key
		var key = int(binary.LittleEndian.Uint64(blob.XXX_unrecognized))

		d.Index.Blobs[key].Groups = append(d.Index.Blobs[key].Groups, info)
		d.Mutex.Unlock()

	}

	return nil
}

// should be concurrency safe
func (d *Decoder) blobData(blob *OSMPBF.Blob) (*OSMPBF.PrimitiveBlock, error) {
	buf := make([]byte, blob.GetRawSize())
	switch {
	case blob.Raw != nil:
		buf = blob.Raw
	case blob.ZlibData != nil:
		r, err := zlib.NewReader(bytes.NewReader(blob.GetZlibData()))
		if err != nil {
			return nil, err
		}
		defer r.Close()

		n, err := io.ReadFull(r, buf)
		if err != nil {
			return nil, err
		}
		if n != int(blob.GetRawSize()) {
			return nil, fmt.Errorf("expected %v bytes, read %v", blob.GetRawSize(), n)
		}
	default:
		return nil, fmt.Errorf("found block with unknown data")
	}
	var primitiveBlock = OSMPBF.PrimitiveBlock{}
	err := proto.Unmarshal(buf, &primitiveBlock)
	return &primitiveBlock, err
}
