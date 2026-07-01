package main

import (
	"fmt"
	"golang-postrgres-crud/db"
	"golang-postrgres-crud/handlers"

	_ "golang-postrgres-crud/docs" // this is importent other wise 500 doc.json

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Golang PostgresSql CRUD API
// @version 1.0
// @description A structural segmented CRUD API using Gin, PostgreSQL, and Swagger with Joins.
// @host localhost:8080
// @BasePath /
func main() {
	db.InitDB()
	defer db.Pool.Close()
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Maps specific routing definitions loading dynamic documentation systems

	r.POST("/employees", handlers.CreateEmployee)
	r.GET("/employees", handlers.GetAllEmployees)
	r.GET("/employees/:id", handlers.GetEmployee)
	r.PUT("/employees/:id", handlers.UpdateEmployee)
	r.DELETE("/employees/:id", handlers.DeleteEmployee)

	r.POST("/department", handlers.CreateDepartment)
	r.GET("/department", handlers.GetAllDepartments)
	r.GET("/department/:id", handlers.GetDepartment)
	r.PUT("/department/:id", handlers.UpdateDepartment)
	r.DELETE("/department/:id", handlers.DeleteDepartment)

	r.POST("/employeedetail", handlers.CreateEmployeeDetail)

	fmt.Println("Structured server running smoothly on port 8080...")
	r.Run(":8080")
}
