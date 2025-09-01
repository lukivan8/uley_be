package router

import (
	"os"
	"strconv"
	"strings"
	"time"

	"uley_be/services"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterRoutes(app core.App) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/collections/v2/items", func(e *core.RequestEvent) error {
			q := e.Request.URL.Query()

			limit, _ := strconv.Atoi(q.Get("limit"))
			offset, _ := strconv.Atoi(q.Get("offset"))
			sort := q.Get("sort")

			var maxP *float64
			if v := strings.TrimSpace(q.Get("max_price")); v != "" {
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					maxP = &f
				}
			}

			items, err := services.ListItems(e.App, services.ItemsFilter{
				MaxPrice:   maxP,
				Location:   q.Get("location"),
				Search:     q.Get("search"),
				CategoryID: q.Get("category_id"),
				Limit:      limit,
				Offset:     offset,
				Sort:       sort,
			})
			if err != nil {
				return e.JSON(500, map[string]any{"error": err.Error()})
			}

			return e.JSON(200, items)
		})

		se.Router.GET("/api/collections/v2/items/{id}", func(e *core.RequestEvent) error {
			id := e.Request.PathValue("id")
			item, err := services.GetItem(e.App, id)
			if err != nil {
				return e.JSON(500, map[string]any{"error": err.Error()})
			}
			if item == nil {
				return e.JSON(404, map[string]any{"error": "item not found"})
			}
			return e.JSON(200, item)
		})

		se.Router.GET("/api/rents/{item_id}", func(e *core.RequestEvent) error {
			item_id := e.Request.PathValue("item_id")
			if item_id == "" {
				return e.JSON(400, map[string]any{"error": "item_id is required"})
			}

			rent, err := services.GetRentByItemID(e.App, item_id)
			if err != nil {
				return e.JSON(500, map[string]any{"error": err.Error()})
			}
			if rent == nil {
				return e.JSON(404, map[string]any{"error": "rent not found"})
			}

			return e.JSON(200, rent)
		})

		se.Router.POST("/api/rents/{item_id}", func(e *core.RequestEvent) error {
			q := e.Request.URL.Query()
			params := services.RentItemRequest{}

			itemID := e.Request.PathValue("item_id")
			if itemID == "" {
				return e.JSON(400, map[string]any{"error": "item_id is required"})
			}
			params.ItemID = itemID

			renterID := q.Get("renter_id")
			if renterID == "" {
				return e.JSON(400, map[string]any{"error": "renter_id is required"})
			}
			params.RenterID = renterID

			dateStartStr := q.Get("date_start")
			if dateStartStr == "" {
				return e.JSON(400, map[string]any{"error": "date_start is required"})
			}
			dateStart, err := time.Parse(time.RFC3339, dateStartStr)
			if err != nil {
				return e.JSON(400, map[string]any{"error": "invalid date_start format, must be RFC3339"})
			}
			params.DateStart = dateStart

			dateEndStr := q.Get("date_end")
			if dateEndStr == "" {
				return e.JSON(400, map[string]any{"error": "date_end is required"})
			}
			dateEnd, err := time.Parse(time.RFC3339, dateEndStr)
			if err != nil {
				return e.JSON(400, map[string]any{"error": "invalid date_end format, must be RFC3339"})
			}
			params.DateEnd = dateEnd

			rents, err := services.RentItem(e.App, params)
			if err != nil {
				return e.JSON(500, map[string]any{"error": err.Error()})
			}

			return e.JSON(200, map[string]any{"rents": rents})
		})

		se.Router.GET("/api/rents/{item_id}/days", func(e *core.RequestEvent) error {
			itemID := e.Request.PathValue("item_id")
			if itemID == "" {
				return e.JSON(400, map[string]any{"error": "item_id is required"})
			}

			days, err := services.RentedDays(e.App, itemID)
			if err != nil {
				return e.JSON(500, map[string]any{"error": err.Error()})
			}

			return e.JSON(200, days)
		})

		// статика
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		return se.Next()
	})
}
