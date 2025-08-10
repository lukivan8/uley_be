package services

import (
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type ItemsFilter struct {
	MaxPrice   *float64
	Location   string
	Search     string
	CategoryID string
	Limit      int
	Offset     int
	Sort       string
}

type ItemsResponse struct {
	Items []map[string]any `json:"items"`
	Total int              `json:"total"`
}

func ListItems(app core.App, f ItemsFilter) (ItemsResponse, error) {
	parts := []string{}
	params := dbx.Params{}

	if f.MaxPrice != nil {
		parts = append(parts, "price <= {:max_price}")
		params["max_price"] = *f.MaxPrice
	}
	if v := strings.TrimSpace(f.Location); v != "" {
		parts = append(parts, "location ~ {:location}")
		params["location"] = v
	}
	if v := strings.TrimSpace(f.Search); v != "" {
		parts = append(parts, "(title ~ {:search} || description ~ {:search} || tags ~ {:search})")
		params["search"] = v
	}
	if v := strings.TrimSpace(f.CategoryID); v != "" {
		parts = append(parts, "category = {:category}")
		params["category"] = v
	}

	filter := strings.Join(parts, " && ")

	total := 0
	query := app.DB().Select("COUNT(*)").From("items")
	if filter != "" {
		query = query.Where(dbx.NewExp(filter, params))
	}
	err := query.Row(&total)
	if err != nil {
		return ItemsResponse{}, err
	}

	sort := f.Sort
	if sort == "" {
		sort = "-created"
	}

	records, err := app.FindRecordsByFilter("items", filter, sort, f.Limit, f.Offset, params)
	if err != nil {
		return ItemsResponse{}, err
	}

	_ = app.ExpandRecords(records, []string{"category", "author", "photos"}, nil)

	items := make([]map[string]any, len(records))
	for i, r := range records {
		items[i] = r.PublicExport()
	}

	return ItemsResponse{
		Items: items,
		Total: total,
	}, nil
}

func GetItem(app core.App, id string) (map[string]any, error) {
	record, err := app.FindRecordById("items", id)
	if err != nil {
		return nil, err
	}

	if record == nil {
		return nil, nil
	}

	_ = app.ExpandRecords([]*core.Record{record}, []string{"category", "author", "photos"}, nil)
	return record.PublicExport(), nil
}
