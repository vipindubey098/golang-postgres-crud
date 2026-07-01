package handlers

import (
	"context"
	"golang-postrgres-crud/db"
	"golang-postrgres-crud/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Saves an employee into db
// @Tags employees
// @Accept json
// @Produce json
// @Param employee body models.Employee true "Employee JSON payload"
// @Success 201 {object} models.Employee
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employees [post]
func CreateEmployee(c *gin.Context) {
	var emp models.Employee
	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	query := `INSERT INTO employees (name, email, department_id) VALUES ($1, $2, $3) RETURNING ID`
	err := db.Pool.QueryRow(context.Background(), query, emp.Name, emp.Email, emp.DepartmentID).Scan(&emp.ID) // Executes data writes returning identities
	// 	err := db.Pool.QueryRow(
	//     context.Background(),
	//     query,
	//     emp.Name,         // This fills $1
	//     emp.Email,        // This fills $2
	//     emp.DepartmentID  // This fills $3
	// ).Scan(&emp.ID)
	//In Go (especially when using the popular jackc/pgx driver for PostgreSQL), db.Pool represents a Connection Pool.
	// Instead of opening a brand new connection to the database every single time you want to run a query (which is incredibly slow and resource-heavy), a connection pool maintains a "pool" of active, open connections that your application can share.
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, emp) // Sends confirmations holding structured model data instances
}

// GetAllEmployees godoc
// @Summary      Get all employees with departments
// @Description  Retrieves all employees joined with their department names using SQL JOIN
// @Tags         employees
// @Produce      json
// @Success      200 {array} models.EmployeeDetail
// @Failure      500 {object} models.ErrorResponse
// @Router       /employees [get]
func GetAllEmployees(c *gin.Context) {
	query := `
		SELECT e.id, e.name, e.email, d.name FROM employees e INNER JOIN department d on e.department_id = d.id`
	rows, err := db.Pool.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	defer rows.Close() // Schedules connection stream closures handling underlying resources cleanups

	employees := []models.EmployeeDetail{} // Instantiates clear typed slices aggregating read database variables
	for rows.Next() {
		var ed models.EmployeeDetail
		if err := rows.Scan(&ed.ID, &ed.Name, &ed.Email, &ed.DepartmentName); err != nil { // Extracts stream attributes maps to structs
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()}) // Delivers trace information outlining scanning parsing structural failures
			return                                                                           // Terminates broken query data transformation processing actions
		} // Resolves row schema scanner state verification processes
		employees = append(employees, ed)
	}
	c.JSON(http.StatusOK, employees)
}

// GetEmployee godoc
// @Summary      Get an employee by ID
// @Description  Retrieves an individual employee from the database by ID
// @Tags         employees
// @Produce      json
// @Param        id path int true "Employee ID"
// @Success      200 {object} models.Employee
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /employees/{id} [get]
func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	var emp models.Employee

	query := `SELECT id, name, email, department_id FROM employees WHERE id = $1`
	err := db.Pool.QueryRow(context.Background(), query, id).Scan(&emp.ID, &emp.Name, &emp.Email, &emp.DepartmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, emp)
}

// UpdateEmployee godoc
// @Summary      Update an employee
// @Description  Modifies field information of an existing employee by ID
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id path int true "Employee ID"
// @Param        employee body models.Employee true "Updated fields object"
// @Success      200 {object} models.SuccessResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /employees/{id} [put]
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var emp models.Employee

	if err := c.ShouldBindJSON(&emp); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	query := `UPDATE employees SET name=$1, email=$2, department_id=$3 where id = $4`
	_, err := db.Pool.Exec(context.Background(), query, emp.Name, emp.Email, emp.DepartmentID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Employee Updated successfully"})
}

// DeleteEmployee godoc
// @Summary      Delete an employee
// @Description  Deletes an employee record from the system by ID
// @Tags         employees
// @Produce      json
// @Param        id path int true "Employee ID"
// @Success      200 {object} models.SuccessResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /employees/{id} [delete]
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM employees WHERE id = $1`
	_, err := db.Pool.Exec(context.Background(), query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Employee deleted successfully"})
}

// CreateDepartment godoc
// @Summary Create a new Department
// @Description Saves an Department into db
// @Tags departments
// @Accept json
// @Produce json
// @Param department body models.Department true "Department JSON payload"
// @Success 201 {object} models.Department
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /department [post]
func CreateDepartment(c *gin.Context) {
	var department models.Department

	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	query := `INSERT INTO department (name) VALUES ($1) RETURNING ID`
	err := db.Pool.QueryRow(context.Background(), query, department.Name).Scan(&department.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, department)
}

// GetAllDepartments godoc
// @Summary      Get all department
// @Description  Retrieves all departments
// @Tags         departments
// @Produce      json
// @Success      200 {array} models.Department
// @Failure      500 {object} models.ErrorResponse
// @Router       /department [get]
func GetAllDepartments(c *gin.Context) {
	query := `SELECT id, name from department`
	rows, err := db.Pool.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var d models.Department

		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		departments = append(departments, d)
	}
	c.JSON(http.StatusOK, departments)
}

// GetDepartment godoc
// @Summary      Get an department by ID
// @Description  Retrieves an individual department from the database by ID
// @Tags         departments
// @Produce      json
// @Param        id path int true "Department ID"
// @Success      200 {object} models.Department
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /department/{id} [get]
func GetDepartment(c *gin.Context) {
	id := c.Param("id")
	var dep models.Department

	query := `SELECT id, name FROM department WHERE id = $1`
	err := db.Pool.QueryRow(context.Background(), query, id).Scan(&dep.ID, &dep.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dep)
}

// UpdateDepartment godoc
// @Summary      Update an department
// @Description  Modifies field information of an existing department by ID
// @Tags         departments
// @Accept       json
// @Produce      json
// @Param        id path int true "Department ID"
// @Param        department body models.Department true "Updated fields object"
// @Success      200 {object} models.SuccessResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /department/{id} [put]
func UpdateDepartment(c *gin.Context) {
	id := c.Param("id")
	var dep models.Employee

	if err := c.ShouldBindJSON(&dep); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	query := `UPDATE department SET name=$1 where id = $2`
	_, err := db.Pool.Exec(context.Background(), query, dep.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Department Updated successfully"})
}

// DeleteDepartment godoc
// @Summary      Delete an department
// @Description  Deletes an department record from the system by ID
// @Tags         departments
// @Produce      json
// @Param        id path int true "Department ID"
// @Success      200 {object} models.SuccessResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /department/{id} [delete]
func DeleteDepartment(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM department WHERE id = $1`
	_, err := db.Pool.Exec(context.Background(), query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Department deleted successfully"})
}

// CreateEmployeeDetail godoc
// @Summary Create a new EmployeeDetail
// @Description Saves an EMployeeDetails into db
// @Tags employeedetails
// @Accept json
// @Produce json
// @Param employeedetail body models.EmployeeDetail true "EmployeeDetail JSON payload"
// @Success 201 {object} models.EmployeeDetail
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employeedetail [post]
func CreateEmployeeDetail(c *gin.Context) {
	var EmployeeDetail models.EmployeeDetail

	if err := c.ShouldBindJSON(&EmployeeDetail); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	query := `INSERT INTO employee_details (name, email, department_name) VALUES ($1, $2, $3) RETURNING ID`
	err := db.Pool.QueryRow(context.Background(), query, EmployeeDetail.Name, EmployeeDetail.Email, EmployeeDetail.DepartmentName).Scan(&EmployeeDetail.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, EmployeeDetail)
}
