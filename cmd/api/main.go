package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"saas_identidad/ent"
	_ "saas_identidad/ent/runtime" // <--- ESTA LÍNEA ES VITAL
	"saas_identidad/handler"
	"saas_identidad/pkg/validation"
	"saas_identidad/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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

	app:=echo.New()

	// create services
	SInvitation:=services.NewInvitationServices(client)
	SUser:=services.NewUserServices(client)
	// create handler
	HInvitation:=handler.NewInvitationHandler(SInvitation,validation.Validator)
	HUser:=handler.NewUserHandler(SUser,validation.Validator)
	api:=app.Group("/tenant/v1")
	// create rest
	RInvitation:=api.Group("/invitation")
	{
		RInvitation.POST("",HInvitation.InvitationDeveloper)
	}
	RUser:=api.Group("/user")
	{
		RUser.POST("",HUser.CreateUser)
	}

	app.Logger.Fatal(app.Start(":8080"))
}