package reader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	defaultDataChannelSize = 10
)

var (
	ErrFirstJSONTokenIncorrect = errors.New("First token in file is not `[`")
)

type JsonDecoderConfig struct {
	DataChannelBufferSize uint
	LogStuffToStdout      bool
}

func getDefaultJsonDecoderConfig() JsonDecoderConfig {
	return JsonDecoderConfig{
		DataChannelBufferSize: defaultDataChannelSize,
		LogStuffToStdout:      true,
	}
}

type JsonDecoderOption func(cfg *JsonDecoderConfig)

func WithCustomDataChannelSize(size uint) JsonDecoderOption {
	return func(cfg *JsonDecoderConfig) {
		cfg.DataChannelBufferSize = size
	}
}

func WithSuppressedLogging() JsonDecoderOption {
	return func(cfg *JsonDecoderConfig) {
		cfg.LogStuffToStdout = false
	}
}

// StartJsonDecoder creates goroutine and launches json.Decoder on file (e.g os.Open(filename))
// Creates and returns chan *pb.Data
// Decoder will close returned channel when it will hit the end of file
func StartJsonDecoder[T any](filename string, opts ...JsonDecoderOption) (chan *T, error) {
	cfg := getDefaultJsonDecoderConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}

	decoder := json.NewDecoder(file)
	tok, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("decoder.Token: %w", err)
	}
	if tok != json.Delim('[') {
		return nil, fmt.Errorf("first token in file %s is not a [ but %v", filename, tok)
	}

	dataChannel := make(chan *T, cfg.DataChannelBufferSize)
	go func() {
		defer file.Close()
		defer close(dataChannel)
		for decoder.More() {
			data := new(T)
			err := decoder.Decode(data)
			if cfg.LogStuffToStdout {
				log.Printf("Got data %v\n", data)
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					if cfg.LogStuffToStdout {
						log.Printf("File ended!")
					}
					break
				}
				log.Printf("ERROR: decoder.Decode: %v", err)
				time.Sleep(time.Second)
				continue
			}
			dataChannel <- data
		}
	}()
	return dataChannel, nil
}
