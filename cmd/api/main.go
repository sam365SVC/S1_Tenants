package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"saas_identidad/ent"
	"saas_identidad/pkg/validation"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	validation.InitValidator()
	err:=godotenv.Load()
	if err!=nil {
		log.Fatalf("Error al cargar el archivo .env: %v",err)
	}

	host:=os.Getenv("PGHOST")
	user:=os.Getenv("PGUSER")
	port:=os.Getenv("PGPORT")
	dbname:=os.Getenv("PGDATABASE")
	dbpassword:=os.Getenv("PGPASSWORD")

	dsn:=fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
						host,port,user,dbpassword,dbname)
	log.Println("Conectandose a Azure PostgreSQL ...")
	
	client,err:=ent.Open("postgres",dsn)
	if err!=nil {
		log.Fatalf("Error al crear schema de DB: %v",err)
	}
	if err:=client.Schema.Create(context.Background());err!=nil {
		log.Fatalf("Error al crear Schemas de DB: %v",err)
	}
	defer client.Close()
}