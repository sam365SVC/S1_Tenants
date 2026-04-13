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
	"saas_identidad/security"

	"saas_identidad/services"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	validation.InitValidator()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	host := os.Getenv("PGHOST")
	user := os.Getenv("PGUSER")
	port := os.Getenv("PGPORT")
	dbname := os.Getenv("PGDATABASE")
	dbpassword := os.Getenv("PGPASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, dbpassword, dbname)
	log.Println("Conectandose a suparbase ...")

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error al crear schema de DB: %v", err)
	}
	if err := client.Schema.Create(
		context.Background(),
		schema.WithDropColumn(true), // <-- Esto le da permiso de BORRAR columnas
		schema.WithDropIndex(true)); err != nil {
		// <-- Esto le da permiso de BORRAR índices); err != nil {
		log.Fatalf("Error al crear Schemas de DB: %v", err)
	}
	defer client.Close()

	app := echo.New()

	// create services
	SLogin := services.NewLoginServices(client)
	SInvitation := services.NewInvitationServices(client)
	SUser := services.NewUserServices(client)
	SPlan := services.NewPlanServices(client)
	SOrganization := services.NewOrganizationServices(client)
	// create handler
	HLogin := handler.NewAuthHandler(SLogin, validation.Validator)
	HInvitation := handler.NewInvitationHandler(SInvitation, validation.Validator)
	HUser := handler.NewUserHandler(SUser, validation.Validator)
	HPlan := handler.NewPlanHandler(SPlan, validation.Validator)
	HOrganization := handler.NewOrganizationHandler(SOrganization, validation.Validator)
	// validation jwt
	jwtAuth := security.JWTMiddleware(os.Getenv("JWT_KEY"))
	// create group api
	api := app.Group("/tenant/v1")
	// create rest
	RInvitation := api.Group("/invitation")
	{
		RInvitation.POST("", HInvitation.InvitationUserOrAdmin)
	}
	RInvitationDeveloper := RInvitation.Group("/developer", jwtAuth, security.RequireRoles("DEVELOPER"))
	{
		RInvitationDeveloper.POST("", HInvitation.InvitationDeveloper)
	}

	RUser := api.Group("/user")
	{
		RUser.POST("", HUser.CreateUser)
	}
	ROrganization := api.Group("/organization", jwtAuth, security.RequireRoles("ADMIN"))
	{
		ROrganization.POST("", HOrganization.CreateOrganization)
	}
	RPlan := api.Group("/plan", jwtAuth, security.RequireRoles("DEVELOPER"))
	{
		RPlan.POST("", HPlan.CreatePlan)
	}

	RLogin := api.Group("/login")
	{
		RLogin.POST("", HLogin.Login)
	}

	app.Logger.Fatal(app.Start(":8080"))
}
