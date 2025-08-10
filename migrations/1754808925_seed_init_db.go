package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// collections
		usersCol, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		categoriesCol, err := app.FindCollectionByNameOrId("categories")
		if err != nil {
			return err
		}
		itemsCol, err := app.FindCollectionByNameOrId("items")
		if err != nil {
			return err
		}
		itemPhotosCol, err := app.FindCollectionByNameOrId("item_photos")
		if err != nil {
			return err
		}
		rentsCol, err := app.FindCollectionByNameOrId("rents")
		if err != nil {
			return err
		}
		favoritesCol, err := app.FindCollectionByNameOrId("favorite_items")
		if err != nil {
			return err
		}

		// seed users (auth records)
		type userSeed struct {
			Email     string
			Password  string
			FirstName string
			LastName  string
			Avatar    string
			Identity  string
			Phone     string
			Verified  bool
		}

		userSeeds := []userSeed{
			{
				Email:     "alina@uley.kz",
				Password:  "AlinaPass123",
				FirstName: "Алина",
				LastName:  "Иванова",
				Avatar:    "https://picsum.photos/seed/alina/200/200",
				Identity:  "990101123456",
				Phone:     "7012345678",
				Verified:  true,
			},
			{
				Email:     "timur@uley.kz",
				Password:  "TimurPass123",
				FirstName: "Тимур",
				LastName:  "Сапаров",
				Avatar:    "https://picsum.photos/seed/timur/200/200",
				Identity:  "920202654321",
				Phone:     "7023456789",
				Verified:  true,
			},
			{
				Email:     "dina@uley.kz",
				Password:  "DinaPass123",
				FirstName: "Дина",
				LastName:  "Нурхан",
				Avatar:    "https://picsum.photos/seed/dina/200/200",
				Identity:  "010303987654",
				Phone:     "7055555555",
				Verified:  true,
			},
		}

		createdUsers := make(map[string]*core.Record)
		for _, u := range userSeeds {
			// idempotent: skip if already exists
			if existing, _ := app.FindFirstRecordByData(usersCol.Id, "email", u.Email); existing != nil {
				createdUsers[u.Email] = existing
				continue
			}
			rec := core.NewRecord(usersCol)
			rec.Set("email", u.Email)
			rec.Set("password", u.Password)
			rec.Set("first_name", u.FirstName)
			rec.Set("last_name", u.LastName)
			rec.Set("avatar", u.Avatar)
			rec.Set("identity", u.Identity)
			rec.Set("phone", u.Phone)
			rec.Set("verified", u.Verified)
			if err := app.Save(rec); err != nil {
				return err
			}
			createdUsers[u.Email] = rec
		}

		// seed categories
		categoryNames := []string{
			"Инструменты",
			"Электроника",
			"Спорт",
			"Кемпинг",
		}
		createdCategories := make(map[string]*core.Record)
		for _, name := range categoryNames {
			if existing, _ := app.FindFirstRecordByData(categoriesCol.Id, "name", name); existing != nil {
				createdCategories[name] = existing
				continue
			}
			rec := core.NewRecord(categoriesCol)
			rec.Set("name", name)
			if err := app.Save(rec); err != nil {
				return err
			}
			createdCategories[name] = rec
		}

		// helper to fetch category/user ids
		cat := func(name string) string { return createdCategories[name].Id }
		user := func(email string) string { return createdUsers[email].Id }

		// seed items
		type itemSeed struct {
			Title       string
			Price       float64
			Description string
			Location    string
			Tags        string
			Category    string
			AuthorEmail string
			PhotoSeeds  []string
		}

		itemSeeds := []itemSeed{
			{
				Title:       "Перфоратор Bosch GBH 2-26",
				Price:       8000,
				Description: "Надёжный перфоратор для ремонта. Подходит для бетона и кирпича.",
				Location:    "Алматы, Бостандыкский район",
				Tags:        "инструменты, ремонт",
				Category:    "Инструменты",
				AuthorEmail: "alina@uley.kz",
				PhotoSeeds:  []string{"bosch1", "bosch2"},
			},
			{
				Title:       "Палатка 3-местная NatureHike",
				Price:       6000,
				Description: "Лёгкая и прочная палатка для выходных на природе.",
				Location:    "Алматы, Медеуский район",
				Tags:        "кемпинг, отдых",
				Category:    "Кемпинг",
				AuthorEmail: "timur@uley.kz",
				PhotoSeeds:  []string{"tent1", "tent2"},
			},
			{
				Title:       "Велосипед горный Trek Marlin 7",
				Price:       10000,
				Description: "Отлично подходит для прогулок и трейлов. Настроен и готов к поездке.",
				Location:    "Астана, Алматы район",
				Tags:        "спорт, велосипед",
				Category:    "Спорт",
				AuthorEmail: "dina@uley.kz",
				PhotoSeeds:  []string{"bike1", "bike2"},
			},
			{
				Title:       "Проектор Xiaomi Mi Smart",
				Price:       15000,
				Description: "Яркий проектор для фильмов и презентаций. Поддержка Wi‑Fi.",
				Location:    "Астана, Есильский район",
				Tags:        "электроника, кино",
				Category:    "Электроника",
				AuthorEmail: "timur@uley.kz",
				PhotoSeeds:  []string{"proj1", "proj2"},
			},
		}

		createdItems := make(map[string]*core.Record)
		for _, it := range itemSeeds {
			if existing, _ := app.FindFirstRecordByData(itemsCol.Id, "title", it.Title); existing != nil {
				createdItems[it.Title] = existing
				continue
			}
			rec := core.NewRecord(itemsCol)
			rec.Set("price", it.Price)
			rec.Set("title", it.Title)
			rec.Set("description", it.Description)
			rec.Set("has_photos", true)
			rec.Set("location", it.Location)
			rec.Set("tags", it.Tags)
			rec.Set("category", cat(it.Category))
			rec.Set("author", user(it.AuthorEmail))
			if err := app.Save(rec); err != nil {
				return err
			}
			createdItems[it.Title] = rec

			// photos
			for _, seed := range it.PhotoSeeds {
				url := "https://picsum.photos/seed/uley-" + seed + "/800/600"
				if existingPhoto, _ := app.FindFirstRecordByData(itemPhotosCol.Id, "url", url); existingPhoto != nil {
					continue
				}
				photo := core.NewRecord(itemPhotosCol)
				photo.Set("item", rec.Id)
				photo.Set("url", url)
				if err := app.Save(photo); err != nil {
					return err
				}
			}
		}

		// one rent example for projector (renter Alina)
		if projector := createdItems["Проектор Xiaomi Mi Smart"]; projector != nil {
			// idempotent on exact dates
			if existing, _ := app.FindFirstRecordByData(rentsCol.Id, "item", projector.Id); existing == nil {
				rent := core.NewRecord(rentsCol)
				rent.Set("item", projector.Id)
				rent.Set("renter", user("alina@uley.kz"))
				rent.Set("date_start", "2025-08-12 00:00:00Z")
				rent.Set("date_end", "2025-08-15 00:00:00Z")
				if err := app.Save(rent); err != nil {
					return err
				}
			}
		}

		// one favorite: Dina favorites Bosch
		if bosch := createdItems["Перфоратор Bosch GBH 2-26"]; bosch != nil {
			// prevent duplicates by (item,user)
			records, err := app.FindAllRecords(favoritesCol.Id, dbx.HashExp{"item": bosch.Id, "user": user("dina@uley.kz")})
			if err != nil {
				return err
			}
			if len(records) == 0 {
				fav := core.NewRecord(favoritesCol)
				fav.Set("item", bosch.Id)
				fav.Set("user", user("dina@uley.kz"))
				if err := app.Save(fav); err != nil {
					return err
				}
			}
		}

		return nil
	}, func(app core.App) error {
		// collections
		usersCol, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		categoriesCol, err := app.FindCollectionByNameOrId("categories")
		if err != nil {
			return err
		}
		itemsCol, err := app.FindCollectionByNameOrId("items")
		if err != nil {
			return err
		}
		itemPhotosCol, err := app.FindCollectionByNameOrId("item_photos")
		if err != nil {
			return err
		}
		rentsCol, err := app.FindCollectionByNameOrId("rents")
		if err != nil {
			return err
		}
		favoritesCol, err := app.FindCollectionByNameOrId("favorite_items")
		if err != nil {
			return err
		}

		// delete favorites and rents for seeded items
		itemTitles := []string{
			"Перфоратор Bosch GBH 2-26",
			"Палатка 3-местная NatureHike",
			"Велосипед горный Trek Marlin 7",
			"Проектор Xiaomi Mi Smart",
		}
		for _, title := range itemTitles {
			if itemRec, _ := app.FindFirstRecordByData(itemsCol.Id, "title", title); itemRec != nil {
				// favorites referencing the item
				favs, err := app.FindAllRecords(favoritesCol.Id, dbx.HashExp{"item": itemRec.Id})
				if err != nil {
					return err
				}
				for _, f := range favs {
					if err := app.Delete(f); err != nil {
						return err
					}
				}
				// rents referencing the item
				rents, err := app.FindAllRecords(rentsCol.Id, dbx.HashExp{"item": itemRec.Id})
				if err != nil {
					return err
				}
				for _, r := range rents {
					if err := app.Delete(r); err != nil {
						return err
					}
				}
				// photos referencing the item
				photos, err := app.FindAllRecords(itemPhotosCol.Id, dbx.HashExp{"item": itemRec.Id})
				if err != nil {
					return err
				}
				for _, p := range photos {
					if err := app.Delete(p); err != nil {
						return err
					}
				}
			}
		}

		// delete items
		for _, title := range itemTitles {
			if rec, _ := app.FindFirstRecordByData(itemsCol.Id, "title", title); rec != nil {
				if err := app.Delete(rec); err != nil {
					return err
				}
			}
		}

		// delete categories
		categoryNames := []string{"Инструменты", "Электроника", "Спорт", "Кемпинг"}
		for _, name := range categoryNames {
			if rec, _ := app.FindFirstRecordByData(categoriesCol.Id, "name", name); rec != nil {
				if err := app.Delete(rec); err != nil {
					return err
				}
			}
		}

		// delete users
		userEmails := []string{"alina@uley.kz", "timur@uley.kz", "dina@uley.kz"}
		for _, email := range userEmails {
			if rec, _ := app.FindFirstRecordByData(usersCol.Id, "email", email); rec != nil {
				if err := app.Delete(rec); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
