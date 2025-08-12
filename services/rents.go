package services

import (
	"fmt"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type RentItemRequest struct {
	ItemID string
	RenterID string
    DateStart time.Time
    DateEnd   time.Time
}

func GetRentByItemID(app core.App, itemID string) ([]map[string]any, error) {
    records, err := app.FindRecordsByFilter(
        "rents",
        "item = {:item_id}",
        "",
        0,
        0,
        dbx.Params{"item_id": itemID},
    )
    if err != nil {
        return nil, err
    }

    _ = app.ExpandRecords(records, []string{"item"}, nil)

    rents := make([]map[string]any, 0, len(records))
    for _, rec := range records {
        rents = append(rents, rec.PublicExport())
    }

    return rents, nil
}

func RentItem(app core.App, req RentItemRequest) ([]map[string]any, error) {
    collection, err := app.FindCollectionByNameOrId("rents")
    if err != nil {
        return nil, err
    }

    rentItems, err := GetRentByItemID(app, req.ItemID)
    if err != nil {
        return nil, err
    }


    for _, rentMap := range rentItems {
        var startDate, endDate time.Time

        switch v := rentMap["date_start"].(type) {
        case string:
            parsed, err := time.Parse(time.RFC3339, v)
            if err != nil {
                return nil, fmt.Errorf("invalid start date in existing rent: %w", err)
            }
            startDate = parsed
        case time.Time:
            startDate = v
        case types.DateTime:
            startDate = v.Time()
        default:
            return nil, fmt.Errorf("invalid start date type: %T", v)
        }

        switch v := rentMap["date_end"].(type) {
        case string:
            parsed, err := time.Parse(time.RFC3339, v)
            if err != nil {
                return nil, fmt.Errorf("invalid end date in existing rent: %w", err)
            }
            endDate = parsed
        case time.Time:
            endDate = v
        case types.DateTime:
            endDate = v.Time()
        default:
            return nil, fmt.Errorf("invalid end date type: %T", v)
        }

        if req.DateStart.Before(endDate) && req.DateEnd.After(startDate) {
            return nil, fmt.Errorf("item is already rented during the requested period")
        }
    }

    record := core.NewRecord(collection)
    record.Set("item", req.ItemID)
    record.Set("renter", req.RenterID)
    record.Set("date_start", req.DateStart.Format(time.RFC3339))
    record.Set("date_end", req.DateEnd.Format(time.RFC3339))

    err = app.Save(record)
    if err != nil {
        return nil, err
    }

    return []map[string]any{
        {
            "item":       req.ItemID,
            "renter":     req.RenterID,
            "date_start": req.DateStart,
            "date_end":   req.DateEnd,
        },
    }, nil
}
