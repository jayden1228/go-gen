package main

import (
	"go-gen/internal/pkg/database"
	"go-gen/internal/service"
	"log"
	"os"
	"os/signal"

	"github.com/gogf/gf/frame/g"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// init mysql
	database.SetUp()
	ts, err := service.GenCustomTableDesc("damns")
	if err != nil {
		panic(err)
	}

	g.Dump(ts)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("quite service ...")
}
