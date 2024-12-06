package util_test

import (
	"graphql_api/util"
	"testing"
	pb "worker"
)

// Test ApplyArgumentsToPbDataFilter
func TestApplyArgumentsToPbDataFilter(t *testing.T) {
	// Test case 1: All arguments are nil
	req := util.ApplyArgumentsToPbDataFilter(nil, nil, nil)
	if req.Limit != 10 {
		t.Errorf("Expected Limit to be 10, got %d", req.Limit)
	}

	// Test case 2: Title is provided
	title := "test-title"
	req = util.ApplyArgumentsToPbDataFilter(&title, nil, nil)
	if req.Title != title {
		t.Errorf("Expected Title to be %s, got %s", title, req.Title)
	}

	// Test case 3: PageToken is provided
	pageToken := "test-page-token"
	req = util.ApplyArgumentsToPbDataFilter(nil, &pageToken, nil)
	if req.PageToken != pageToken {
		t.Errorf("Expected PageToken to be %s, got %s", pageToken, req.PageToken)
	}

	// Test case 4: Limit is provided
	limit := 20
	req = util.ApplyArgumentsToPbDataFilter(nil, nil, &limit)
	if req.Limit != uint64(limit) {
		t.Errorf("Expected Limit to be %d, got %d", limit, req.Limit)
	}
}

// Test ConvertPbDataToModelData
func TestConvertPbDataToModelData(t *testing.T) {
	// Prepare test input
	pbData := &pb.Data{
		XId: "test-id",
		Categories: &pb.Categories{
			Subcategory: "test-subcategory",
		},
		Title: &pb.MultiLanguageTitle{
			Ru: "Тестовый заголовок",
			Ro: "Test title",
		},
		Type:   "test-type",
		Posted: 1234567890,
	}

	// Call the function
	modelData := util.ConvertPbDataToModelData(pbData)

	// Verify the result
	if modelData.ID != pbData.XId {
		t.Errorf("Expected ID to be %s, got %s", pbData.XId, modelData.ID)
	}
	if modelData.Categories == nil {
		t.Errorf("modelData.Categories is nil!")
	}
	if modelData.Categories.Subcategory != pbData.Categories.Subcategory {
		t.Errorf("Expected Subcategory to be %s, got %s", pbData.Categories.Subcategory, modelData.Categories.Subcategory)
	}
	if modelData.Title == nil {
		t.Errorf("modelData.Title is nil!")
	}
	if modelData.Title.Ru != pbData.Title.Ru {
		t.Errorf("Expected Title.Ru to be %s, got %s", pbData.Title.Ru, modelData.Title.Ru)
	}
	if modelData.Title.Ro != pbData.Title.Ro {
		t.Errorf("Expected Title.Ro to be %s, got %s", pbData.Title.Ro, modelData.Title.Ro)
	}
	if modelData.Type != pbData.Type {
		t.Errorf("Expected Type to be %s, got %s", pbData.Type, modelData.Type)
	}
	if modelData.Posted != pbData.Posted {
		t.Errorf("Expected Posted to be %f, got %f", pbData.Posted, modelData.Posted)
	}
}

// Test ApplyArgumentsToPbAggregateSearch
func TestApplyArgumentsToPbAggregateSearch(t *testing.T) {
	subcategory := "test-subcategory"
	req := util.ApplyArgumentsToPbAggregateSearch(subcategory)

	// Verify the result
	if req.Categories == nil {
		t.Errorf("req.Categories is nil!")
	}
	if req.Categories.Subcategory != subcategory {
		t.Errorf("Expected Subcategory to be %s, got %s", subcategory, req.Categories.Subcategory)
	}
}

// Test ConvertPbAggregatedDataToModelAggregatedCategory
func TestConvertPbAggregatedDataToModelAggregatedCategory(t *testing.T) {
	pbAggregatedData := &pb.AggregatedData{
		Categories: &pb.Categories{
			Subcategory: "test-subcategory",
		},
		Count: 42,
	}
	modelAggregatedCategory := util.ConvertPbAggregatedDataToModelAggregatedCategory(pbAggregatedData)

	// Verify the result
	if modelAggregatedCategory.Category == nil {
		t.Errorf("modelAggregatedCategory.Category is nil!")
	}
	if modelAggregatedCategory.Category.Subcategory != pbAggregatedData.Categories.Subcategory {
		t.Errorf("Expected Subcategory to be %s, got %s", pbAggregatedData.Categories.Subcategory, modelAggregatedCategory.Category.Subcategory)
	}
	if *modelAggregatedCategory.Count != int(pbAggregatedData.Count) {
		t.Errorf("Expected Count to be %d, got %d", pbAggregatedData.Count, *modelAggregatedCategory.Count)
	}
}
