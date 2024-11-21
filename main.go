package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
)

const file string = "employees.db"

var DB *sql.DB

//	var employees = []employee{
//		{employee_id: 1, department: "Engineering", job_title: "Senior Enginer"},
//		{employee_id: 2, department: "Engineering", job_title: "Super Senior Enginer"},
//		{employee_id: 3, department: "Sales", job_title: "Head of Sales"},
//		{employee_id: 4, department: "Support", job_title: "Tech Support"},
//		{employee_id: 5, department: "Engineering", job_title: "Junior Enginer"},
//		{employee_id: 6, department: "Sales", job_title: "Sales Rep"},
//		{employee_id: 7, department: "Marketing", job_title: "Senior Marketer"},
//	}

// get all employees info
func getEmployees(ctx *gin.Context) {
	employees, err := getAllEmployees()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.IndentedJSON(http.StatusOK, employees)
}

// get Unique employee info by id
func getEmployeeById(ctx *gin.Context) {
	id := ctx.Param("id")

	employee, err := GetEmployeeById(id)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Employee with id: " + id + " does not exist."})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": employee})
}

// update Employee based on id
func updateEmployee(ctx *gin.Context) {
	var json Employee
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	EmployeeId := ctx.Param("id")

	success, err := UpdateEmployee(json, EmployeeId)
	if success {
		ctx.JSON(http.StatusOK, gin.H{"message": "Updated Employee with id: " + EmployeeId})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// add new Employee
func addEmployee(ctx *gin.Context) {

	var json Employee

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, _ := GetEmployeeById(strconv.Itoa(json.Id))
	if employee != (Employee{}) {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"message": "Employee with id " + strconv.Itoa(json.Id) + " already exists."})
		return
	}

	success, err := AddEmployee(json)

	if success {
		ctx.JSON(http.StatusOK, gin.H{"message": "Employee id: " + strconv.Itoa(json.Id) + " added."})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// delete Employee based on id
func deleteEmployee(ctx *gin.Context) {
	id := ctx.Param("id")

	success, err := DeleteEmployee(id)

	if success {
		ctx.JSON(http.StatusOK, gin.H{"message": "Employee id: " + id + " deleted successfully."})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func main() {
	ConnectDatabase()

	r := gin.Default()

	// API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("employees", getEmployees)
		v1.GET("employee/:id", getEmployeeById)
		v1.POST("employee", addEmployee)
		v1.PUT("employee/:id", updateEmployee)
		v1.DELETE("employee/:id", deleteEmployee)
	}

	port := os.Getenv("PORT")
	if port != "" {
		r.Run("localhost:" + port)
	} else {
		log.Fatal("PORT not specified")
	}
}
