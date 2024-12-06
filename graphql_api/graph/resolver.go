package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"context"
	"fmt"
	"graphql_api/graph/model"
	"graphql_api/grpc"
	"graphql_api/util"
	"log"
)

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type Resolver struct{}

type queryResolver struct{ *Resolver }

// Data is the resolver for the data field.
func (r *queryResolver) Data(ctx context.Context, title *string, pageToken *string, limit *int) (*model.DataPage, error) {
	// Default value for req.Limit is 10
	// so if `limit` is `nil` it req.Limit will be equal to 10
	req := util.ApplyArgumentsToPbDataFilter(title, pageToken, limit)

	// You can actually try retrying connection, you need to handle error like that:
	// if errors.Is(err, grpc.ErrConnectionWasNotInitialized) {
	// and call grpc.SafeGetConnection() again
	g, err := grpc.SafeGetConnection()
	if err != nil {
		return nil, fmt.Errorf("grpc.SafeGetConnection: %w", err)
	}
	defer g.Conn.Close()

	dataChannel, err := g.GetStreamOfDataFromServer(context.Background(), req, grpc.WithStreamChannelSizeOf(8))
	if err != nil {
		return nil, fmt.Errorf("conn.GetStreamOfDataFromServer: %w", err)
	}

	var dataArray []*model.Data
	var nextPageIDData string
	for data := range dataChannel {
		log.Printf("Get from channel %v\n", data)
		dataArray = append(dataArray, util.ConvertPbDataToModelData(data))
		nextPageIDData = fmt.Sprintf("%f", data.Posted)
	}

	var nextPageID *string
	if nextPageIDData != "" {
		nextPageID = &nextPageIDData
	}
	ret := &model.DataPage{
		Data:          dataArray,
		NextPageToken: nextPageID,
	}
	return ret, nil
}

// func (r *queryResolver) AggregateSubcategory(ctx context.Context, subcategory string) ([]*model., error) {
func (r *queryResolver) AggregateSubcategory(ctx context.Context, subcategory string) ([]*model.AggregatedCategory, error) {
	req := util.ApplyArgumentsToPbAggregateSearch(subcategory)

	g, err := grpc.SafeGetConnection()
	if err != nil {
		return nil, fmt.Errorf("grpc.SafeGetConnection: %w", err)
	}
	defer g.Conn.Close()

	resp, err := g.GetAggregatedDataBySubcategory(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("conn.GetStreamOfDataFromServer: %w", err)
	}

	var dataArray []*model.AggregatedCategory
	for _, pbCat := range resp.Data {
		dataArray = append(dataArray, util.ConvertPbAggregatedDataToModelAggregatedCategory(pbCat))
	}
	return dataArray, nil
}
