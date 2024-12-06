package reader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	pb "worker"
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
func StartJsonDecoder(filename string, opts ...JsonDecoderOption) (chan *pb.Data, error) {
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
	// this ensures that file is at least some sort of array
	if delim, ok := tok.(json.Delim); ok {
		if delim != '[' {
			log.Printf("First token in file %s is %v with type of %T\n", filename, tok, tok)
			return nil, ErrFirstJSONTokenIncorrect
		}
	} else {
		log.Printf("WARN: First token in file %s is NOT a delim but it is %v with type %T\n", filename, tok, tok)
	}
	dataChannel := make(chan *pb.Data, cfg.DataChannelBufferSize)
	go func() {
		defer close(dataChannel)
		for {
			data := &pb.Data{}
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
				// As decoder.Decode calls json.Decoder internal function
				// ```
				// // advance tokenstate from a separator state to a value state
				// func (dec *Decoder) tokenPrepareForDecode() error {
				// ```
				// we can get this implicit string error that can not be cought with
				// `errors.Is()` as it is not a variable to compare error with
				// ```
				// return &SyntaxError{"expected comma after array element", dec.InputOffset()}
				// ```
				// Therefore we must check it in this way as it will save us when we encounter
				// `}]` at the end of the file
				if strings.Contains(err.Error(), "expected comma after array element") {
					// we try to get the delim `]` and end the stream
					tok, err := decoder.Token()
					if err != nil {
						log.Printf("decoder.Token: %v", err)
						return
					}
					if delim, ok := tok.(json.Delim); ok {
						if delim != ']' {
							log.Printf("Potentially last token in file %s is %v with type of %T. It should be `]`!\n", filename, tok, tok)
						}
						if cfg.LogStuffToStdout {
							log.Printf("File ended!")
						}
						// exiting without fuss
						return
					} else {
						log.Printf("Potentially last token in file %s is %v with type of %T. It should be `]`!\n", filename, tok, tok)
						return
					}
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
