package grpc

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"worker_client/elasticsearch"

	pb "worker"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type storageServer struct {
	pb.UnimplementedStorageServer
}

func (s *storageServer) GetStreamOfDataFromElasticSearch(req *pb.DataSearch, stream pb.Storage_GetStreamOfDataFromElasticSearchServer) error {
	es8Conn := elasticsearch.GetEs8Connection()
	esResult, err := elasticsearch.SearchDataBy(es8Conn, req, stream.Context())
	if err != nil {
		return fmt.Errorf("elasticsearch.SearchDataBy: %w", err)
	}
	for _, hit := range esResult.Hits.Hits {
		data := &pb.Data{
			XId: hit.Source.Id,
			Title: &pb.MultiLanguageTitle{
				Ro: hit.Source.Title.Ro,
				Ru: hit.Source.Title.Ru,
			},
			Categories: &pb.Categories{
				Subcategory: hit.Source.Categories.Subcategory,
			},
			Type:   hit.Source.Type,
			Posted: hit.Source.Posted,
		}

		log.Printf("Sending %v to GraphQL!\n", data)

		if err := stream.Send(data); err != nil {
			return fmt.Errorf("error sending data to stream: %w", err)
		}

		// Обновляем PageToken для клиента
		if len(hit.Sort) > 0 {
			pageToken, _ := json.Marshal(hit.Sort)
			// Клиенту нужно сохранить это значение для следующего запроса
			log.Printf("Next PageToken: %s\n", pageToken)
		}
	}

	return nil
}

func (s *storageServer) GetAggregatedData(ctx context.Context, req *pb.AggregatedDataSearch) (*pb.AggregatedDataRespone, error) {
	es8Conn := elasticsearch.GetEs8Connection()
	esResult, err := elasticsearch.SearchAggregatedDataBy(es8Conn, req, ctx)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch.SearchAggregatedDataBy: %w", err)
	}

	var data []*pb.AggregatedData
	for _, bucket := range esResult.Aggregations.SubcategoryCounts.Buckets {
		data = append(data, &pb.AggregatedData{
			Categories: &pb.Categories{
				Subcategory: bucket.Key,
			},
			Count: uint64(bucket.DocCount),
		})
	}
	resp := &pb.AggregatedDataRespone{
		Data: data,
	}
	return resp, nil
}

func StartStorageServer() {
	// Parses if someone put -port= arg
	// if not port is equal to 50051
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

	pb.RegisterStorageServer(s, &storageServer{})
	log.Printf("Server is listening at %s", listener.Addr())
	if err := s.Serve(listener); err != nil {
		panic(err)
	}
}
