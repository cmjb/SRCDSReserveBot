package main

import (
	"fmt"
	"github.com/go-pg/migrations/v8"
)

func main() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table guilds...")
		_, err := db.Exec(`CREATE TABLE "public"."guilds" (
									"id" bigint NOT NULL,
									"active" boolean NOT NULL,
									"name" text NOT NULL,
									"created_at" bigint NOT NULL
								) WITH (oids = false);`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table my_table...")
		_, err := db.Exec(`DROP TABLE "public"."guilds"`)
		return err
	})
}
