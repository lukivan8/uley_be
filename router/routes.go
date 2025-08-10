package router

import (
	"os"
	"strconv"
	"strings"

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

		// статика
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		return se.Next()
	})
}
