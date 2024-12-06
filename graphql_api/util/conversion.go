package util

import (
	"graphql_api/graph/model"
	pb "worker"
)

// This file is needed to simplify conversion from one complex struct to another
//

// Default value for limit is 10
// so if `limit` is `nil` it req.Limit will be equal to 10
func ApplyArgumentsToPbDataFilter(title *string, pageToken *string, limit *int) *pb.DataSearch {
	req := &pb.DataSearch{}
	if title != nil {
		req.Title = *title
	}
	if pageToken != nil {
		req.PageToken = *pageToken
	}
	req.Limit = 10
	if limit != nil {
		req.Limit = uint64(*limit)
	}
	return req
}

func ConvertPbDataToModelData(data *pb.Data) *model.Data {
	return &model.Data{
		ID: data.XId,
		Categories: &model.Category{
			Subcategory: data.Categories.Subcategory,
		},
		Title: &model.MultiLanguageTitle{
			Ru: data.Title.Ru,
			Ro: data.Title.Ro,
		},
		Type:   data.Type,
		Posted: data.Posted,
	}
}

func ApplyArgumentsToPbAggregateSearch(subcategory string) *pb.AggregatedDataSearch {
	return &pb.AggregatedDataSearch{
		Categories: &pb.Categories{
			Subcategory: subcategory,
		},
	}
}

func ConvertPbAggregatedDataToModelAggregatedCategory(data *pb.AggregatedData) *model.AggregatedCategory {
	intVal := int(data.Count)
	count := &intVal
	return &model.AggregatedCategory{
		Category: &model.Category{
			Subcategory: data.Categories.Subcategory,
		},
		Count: count,
	}
}
