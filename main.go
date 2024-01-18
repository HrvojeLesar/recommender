package main

import (
	"context"
	"log"

	"github.com/HrvojeLesar/recommender/config"
	"github.com/HrvojeLesar/recommender/db"
	"github.com/HrvojeLesar/recommender/global"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error

	config := config.New()
	globalInstances, cancel := global.New(context.Background(), config)

	globalInstances.Instance().Mongo, err = db.Setup(globalInstances, config)
	if err != nil {
		log.Panicln(err)
	}

	cancel()
}
