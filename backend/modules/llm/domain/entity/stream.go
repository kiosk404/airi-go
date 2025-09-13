package entity

import (
	"github.com/cloudwego/eino/schema"
)

//go:generate mockgen -destination=mocks/stream.go -package=mocks . IStreamReader
type IStreamReader interface {
	Recv() (*Message, error)
}

type StreamReader struct {
	einoReader *schema.StreamReader[*Message]
}

func NewStreamReader(einoReader *schema.StreamReader[*schema.Message]) IStreamReader {
	return &StreamReader{
		einoReader: schema.StreamReaderWithConvert(einoReader, ToDOMessage),
	}
}

func (sr *StreamReader) Recv() (message *Message, err error) {
	return sr.einoReader.Recv()
}
