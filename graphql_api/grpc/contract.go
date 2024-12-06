package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	pb "worker"
)

type streamConfig struct {
	channelSize uint
}

func createDefaultStreamConfig() streamConfig {
	return streamConfig{
		channelSize: 1,
	}
}

type StreamConfigOption func(cfg *streamConfig)

func WithStreamChannelSizeOf(size uint) StreamConfigOption {
	return func(cfg *streamConfig) {
		cfg.channelSize = size
	}
}

func (g *GRPCConnection) GetStreamOfDataFromServer(ctx context.Context, req *pb.DataSearch, streamOpts ...StreamConfigOption) (chan *pb.Data, error) {
	cfg := createDefaultStreamConfig()
	for _, opt := range streamOpts {
		opt(&cfg)
	}
	stream, err := g.Client.GetStreamOfDataFromElasticSearch(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("g.Client.GetStreamOfDataFromElasticSearch: %w", err)
	}
	retChan := make(chan *pb.Data, cfg.channelSize)
	go func() {
		defer close(retChan)
		for {
			data, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Printf("ERROR: stream.Recv: %v\n", err)
				return
			}
			retChan <- data
		}
	}()
	return retChan, nil
}

func (g *GRPCConnection) GetAggregatedDataBySubcategory(ctx context.Context, req *pb.AggregatedDataSearch) (*pb.AggregatedDataRespone, error) {
	resp, err := g.Client.GetAggregatedData(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
