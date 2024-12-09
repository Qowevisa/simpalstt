package grpc

import (
	"flag"
	"fmt"
	"log"
	"net"
	"worker_server/env"
	"worker_server/reader"

	pb "worker"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port = flag.Int("port", 50050, "The server port")
)

type server struct {
	pb.UnimplementedWorkerServer
}

// Implements WorkerServer GetStreamOfData rpc call
func (s *server) GetStreamOfData(req *pb.DataFilter, stream pb.Worker_GetStreamOfDataServer) error {
	filterIds := make(map[string]bool)
	for _, idToFilter := range req.Id {
		filterIds[idToFilter] = true
	}
	defaultJsonFilePath := env.JSONFilepath()
	dataChannel, err := reader.StartJsonDecoder[pb.Data](defaultJsonFilePath)
	if err != nil {
		return fmt.Errorf("reader.LaunchJSONDecoderOn: %w", err)
	}
	for data := range dataChannel {
		if _, filterIt := filterIds[data.XId]; filterIt {
			continue
		}
		if err := stream.Send(data); err != nil {
			return err
		}
	}
	return nil
}

func StartWorkerServer() {
	// Parses if someone put -port= arg
	// if not port is equal to 50050
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()

	// Healh server to answer to `grpc_health_probe` executable for
	// docker-compose healthcheck
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	pb.RegisterWorkerServer(s, &server{})
	log.Printf("Server is listening at %s", listener.Addr())
	if err := s.Serve(listener); err != nil {
		panic(err)
	}
}
