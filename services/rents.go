package services

import (
	"fmt"
	"sort"
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
        "-date_start",
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

func dateOnly(t time.Time) time.Time {
    return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
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
            parsed, _ := time.Parse(time.RFC3339, v)
            startDate = parsed
        case time.Time:
            startDate = v
        case types.DateTime:
            startDate = v.Time()
        }

        switch v := rentMap["date_end"].(type) {
        case string:
            parsed, _ := time.Parse(time.RFC3339, v)
            endDate = parsed
        case time.Time:
            endDate = v
        case types.DateTime:
            endDate = v.Time()
        }

        startDate = dateOnly(startDate)
        endDate = dateOnly(endDate)

        endDateExclusive := endDate.AddDate(0, 0, 1)

        reqStart := dateOnly(req.DateStart)
        reqEnd := dateOnly(req.DateEnd).AddDate(0, 0, 1)

        if reqStart.Before(endDateExclusive) && reqEnd.After(startDate) {
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

func RentedDays(app core.App, itemID string) ([]string, error) {
    rentItems, err := GetRentByItemID(app, itemID)
    if err != nil {
        return nil, err
    }

    var allDays []time.Time

    for _, rentMap := range rentItems {
        var startDate, endDate time.Time

        switch v := rentMap["date_start"].(type) {
        case string:
            startDate, _ = time.Parse(time.RFC3339, v)
        case time.Time:
            startDate = v
        case types.DateTime:
            startDate = v.Time()
        }

        switch v := rentMap["date_end"].(type) {
        case string:
            endDate, _ = time.Parse(time.RFC3339, v)
        case time.Time:
            endDate = v
        case types.DateTime:
            endDate = v.Time()
        }

        startDate = dateOnly(startDate)
        endDate = dateOnly(endDate)

        for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
            allDays = append(allDays, d)
        }
    }

    if len(allDays) == 0 {
        return nil, nil
    }

    sort.Slice(allDays, func(i, j int) bool {
        return allDays[i].Before(allDays[j])
    })

    var result []string
    rangeStart := allDays[0]
    prev := allDays[0]

    for i := 1; i < len(allDays); i++ {
        if allDays[i].Sub(prev) > 24*time.Hour {
            if rangeStart.Equal(prev) {
                result = append(result, rangeStart.Format("02.01"))
            } else {
                result = append(result, fmt.Sprintf("%s-%s",
                    rangeStart.Format("02.01"),
                    prev.Format("02.01")))
            }
            rangeStart = allDays[i]
        }
        prev = allDays[i]
    }

    if rangeStart.Equal(prev) {
        result = append(result, rangeStart.Format("02.01"))
    } else {
        result = append(result, fmt.Sprintf("%s-%s",
            rangeStart.Format("02.01"),
            prev.Format("02.01")))
    }

    return result, nil
}