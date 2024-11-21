package main

import (
	"database/sql"
)

type Employee struct {
	Id         int    `json:"employee_id"`
	Gender     string `json:"gender"`
	Department string `json:"department"`
	Job_title  string `json:"job_title"`
}

func checkData(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}

// getEmployees responds with the list of all employees as JSON.
func getAllEmployees() ([]Employee, error) {
	var id int
	var gender, department, job_title sql.NullString
	rows, err := DB.Query("SELECT * from employees")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	employees := make([]Employee, 0)
	for rows.Next() {
		singleEmployee := Employee{}
		err = rows.Scan(&id, &gender, &department, &job_title)
		singleEmployee.Id = id
		singleEmployee.Gender = checkData(gender)
		singleEmployee.Department = checkData(department)
		singleEmployee.Job_title = checkData(job_title)
		if err != nil {
			return nil, err
		}
		employees = append(employees, singleEmployee)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return employees, err
}

func GetEmployeeById(emp_id string) (Employee, error) {
	var id int
	var gender, department, job_title sql.NullString
	stmt, err := DB.Prepare("SELECT id, gender, department, job_title from employees WHERE id = ?")

	if err != nil {
		return Employee{}, err
	}

	employee := Employee{}

	sqlErr := stmt.QueryRow(emp_id).Scan(&id, &gender, &department, &job_title)
	employee.Id = id
	employee.Gender = checkData(gender)
	employee.Department = checkData(department)
	employee.Job_title = checkData(job_title)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Employee{}, sql.ErrNoRows
		}
		return Employee{}, sqlErr
	}
	return employee, nil
}

func AddEmployee(newEmployee Employee) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO employees (id, gender, department, job_title) VALUES (?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newEmployee.Id, newEmployee.Gender, newEmployee.Department, newEmployee.Job_title)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdateEmployee(employee Employee, id string) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE employees SET id = ?, gender = ?, department = ?, job_title = ? WHERE Id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(employee.Id, employee.Gender, employee.Department, employee.Job_title, id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteEmployee(employeeId string) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from employees where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(employeeId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
