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

func ListItems(app core.App, f ItemsFilter) ([]map[string]any, error) {
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

	sort := f.Sort
	if sort == "" {
		sort = "-created"
	}

	records, err := app.FindRecordsByFilter("items", filter, sort, f.Limit, f.Offset, params)
	if err != nil {
		return nil, err
	}

	_ = app.ExpandRecords(records, []string{"category", "author"}, nil)

	result := make([]map[string]any, len(records))
	for i, r := range records {
		result[i] = r.PublicExport()
	}
	return result, nil
}
