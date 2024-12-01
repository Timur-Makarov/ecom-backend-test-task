package main

import (
	"ecom-backend-test-task/config"
	"log"
	"net/http"
)

func main() {
	hardCodedDSN := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := config.OpenDB(hardCodedDSN)
	if err != nil {
		log.Fatalln(err)
	}

	config.CheckIfShouldMigrate(db)

	apiComponents := config.SetupAPIComponents(db)

	apiComponents.Services.BannerService.RunBannerCounterUpdate()

	http.HandleFunc("/banner/add", apiComponents.Handlers.AddBanner)
	http.HandleFunc("/counter/{bannerID}", apiComponents.Handlers.UpdateBannerCounterStats)
	http.HandleFunc("/stats/{bannerID}", apiComponents.Handlers.GetBannerCounterStats)

	log.Println("Listening on 4000 port")

	err = http.ListenAndServe("127.0.0.1:4000", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
