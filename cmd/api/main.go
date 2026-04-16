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
	STenant := services.NewTenantServices(client)
	SEmployee := services.NewEmployeeServices(client)
	SBranch := services.NewBranchServices(client)
	SUser := services.NewUserServices(client)
	SPlan := services.NewPlanServices(client)
	SOrganization := services.NewOrganizationServices(client)
	// create handler
	HLogin := handler.NewAuthHandler(SLogin, validation.Validator)
	HInvitation := handler.NewInvitationHandler(SInvitation, validation.Validator)
	HTenant := handler.NewTenantHandler(STenant, validation.Validator)
	HEmployee := handler.NewEmployeeHandler(SEmployee, validation.Validator)
	HBranch := handler.NewBranchHandler(SBranch, validation.Validator)
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
	RInvitationJob := RInvitation.Group("/job", jwtAuth, security.RequireRoles("ADMIN"))
	{
		RInvitationJob.POST("", HInvitation.InvitationJob)
	}

	RUser := api.Group("/user")
	{
		RUser.POST("", HUser.CreateUser)
		RUser.GET("/:id", HUser.GetUserId)
		RUser.GET("/page/:page", HUser.AllUser)
		RUser.PUT("/:id", HUser.RemplaceUser)
		RUser.PATCH("/:id", HUser.PatchUser)
	}
	RTenant := api.Group("/tenant")
	{
		RTenant.GET("/page/:page", HTenant.GetPageTenant)
	}
	REmployee := api.Group("/employee")
	{
		REmployee.GET("/page/:page", HEmployee.GetPageEmployee)
		REmployee.PUT("/:id", HEmployee.RemplaceEmployee)
		REmployee.PATCH("/:id", HEmployee.PatchEmployee)
	}
	RBranch := api.Group("/branch")
	{
		RBranch.GET("/page/:page", HBranch.GetPageBranch)
	}
	ROrganization := api.Group("/organization", jwtAuth, security.RequireRoles("ADMIN"))
	{
		ROrganization.POST("", HOrganization.CreateOrganization)
	}
	RPlan := api.Group("/plan")
	{
		RPlan.GET("", HPlan.AllPlan)
	}
	RPlanS := RPlan.Group("", jwtAuth, security.RequireRoles("DEVELOPER"))
	{
		RPlanS.POST("", HPlan.CreatePlan)
		RPlanS.PUT("/:id", HPlan.UpdatePlan)
	}
	RLogin := api.Group("/login")
	{
		RLogin.POST("", HLogin.Login)
	}
	RLoginTenant := RLogin.Group("/tenant", jwtAuth, security.RequireRoles("ADMIN"))
	{
		RLoginTenant.POST("", HLogin.LoginTenant)
	}
	app.Logger.Fatal(app.Start(":8080"))
}
