package main

import (
	mysql "crm-service/client"
	"crm-service/config"
	"crm-service/repository"
	"crm-service/route"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	defer mysql.Disconnect()

	{
		loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		time.Local = loc
	}

	cfg := config.GetConfig()
	client := mysql.GetClient
	repo := repository.New(client)

	// generate DB
	//datastores.Migrate(client(context.TODO()))

	{
		h := route.NewHTTPHandler(repo)
		h.Logger.Fatal(h.Start(fmt.Sprintf(":%s", cfg.Port)))
	}
}
