package parser

import (
	"log"
	"os"

	"github.com/missinglink/gosmparse"
)

// Parser - PBF Parser
type Parser struct {
	file    *os.File
	decoder *gosmparse.Decoder
}

// open - open file path
func (p *Parser) open(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	p.file = file
	p.decoder = gosmparse.NewDecoder(file)
}

// Reset - reset (open+close) file
func (p *Parser) Reset() {
	p.file.Close()
	p.open(p.file.Name())
}

// Parse - execute parser
func (p *Parser) Parse(handler gosmparse.OSMReader) {
	err := p.decoder.Parse(handler, false)
	if err != nil {
		panic(err)
	}
}

// ParseFrom - execute parser, starting from offset
func (p *Parser) ParseFrom(handler gosmparse.OSMReader, offset int64) {
	p.decoder.SeekToOffset(offset)
	err := p.decoder.Parse(handler, true)
	if err != nil {
		panic(err)
	}
}

// ParseBlob - execute parser for a single blob
func (p *Parser) ParseBlob(handler gosmparse.OSMReader, offset int64) {
	err := p.decoder.ParseBlob(handler, offset)
	if err != nil {
		panic(err)
	}
}

// GetDecoder - return decoder object
func (p *Parser) GetDecoder() *gosmparse.Decoder {
	return p.decoder
}

// NewParser - Create a new parser for file at path
func NewParser(path string) *Parser {
	p := &Parser{}
	p.open(path)

	return p
}

// NewParserFromArgs - Create a new parser for file at argv position
func NewParserFromArgs(pos int) *Parser {

	if len(os.Args) < (pos + 1) {
		log.Fatal("Invalid argv position")
		os.Exit(1)
	}

	return NewParser(os.Args[pos])
}
