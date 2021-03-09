package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"code.siemens.com/ozdinc.celikel/backend_master_vlass/api"
	db "code.siemens.com/ozdinc.celikel/backend_master_vlass/internal/db"
	"code.siemens.com/ozdinc.celikel/backend_master_vlass/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	tx := db.NewTx(conn)
	server := api.NewServer(tx)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
