package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

func TestListEmployeesByPagination(t *testing.T) {
	store := NewEmployeeStore()

	for i := 1; i <= 20; i++ {
		employee := Employee{
			ID:       i,
			Name:     "Employee" + strconv.Itoa(i),
			Position: "Position" + strconv.Itoa(i),
			Salary:   float64(i * 1000),
		}
		store.Employees[i] = employee
	}

	tests := []struct {
		name           string
		page           int
		limit          int
		expectedStatus int
	}{
		{"ValidPagination", 1, 10, http.StatusOK},
		{"InvalidPage", -1, 10, http.StatusOK},
		{"InvalidLimit", 1, -1, http.StatusOK},
		{"EmptyStore", 1, 10, http.StatusOK},
		{"SecondPage", 2, 10, http.StatusOK},
		{"OutOfRangePage", 100, 10, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/employees?page=%d&limit=%d", tt.page, tt.limit), nil)
			res := httptest.NewRecorder()
			store.ListEmployees(res, req)

			if res.Code != tt.expectedStatus {
				t.Errorf("ListEmployees() %s returned wrong status code: got %v want %v", tt.name, res.Code, tt.expectedStatus)
			}

			if res.Code == http.StatusOK {
				var employees []Employee
				err := json.Unmarshal(res.Body.Bytes(), &employees)
				if err != nil {
					t.Errorf("Error parsing JSON response: %v", err)
				}

				for i := 1; i < len(employees); i++ {
					if employees[i].ID < employees[i-1].ID {
						t.Errorf("ListEmployees() %s returned employees not sorted by ID in ascending order", tt.name)
						break
					}
				}
			}
		})
	}
}

func TestCreateEmployee(t *testing.T) {
	store := NewEmployeeStore()

	employee := Employee{
		ID:       1,
		Name:     "John Doe",
		Position: "Software Engineer",
		Salary:   50000,
	}

	createEmployeeJSON, _ := json.Marshal(employee)
	reqCreateEmployee, _ := http.NewRequest("POST", "/employee", bytes.NewBuffer(createEmployeeJSON))
	reqCreateEmployee.Header.Set("Content-Type", "application/json")
	resCreateEmployee := httptest.NewRecorder()

	store.CreateEmployee(resCreateEmployee, reqCreateEmployee)
	if resCreateEmployee.Code != http.StatusOK {
		t.Errorf("CreateEmployee() returned wrong status code: got %v want %v", resCreateEmployee.Code, http.StatusOK)
	}
}

func TestListEmployees(t *testing.T) {
	store := NewEmployeeStore()

	reqListEmployees, _ := http.NewRequest("GET", "/employees", nil)
	resListEmployees := httptest.NewRecorder()

	store.ListEmployees(resListEmployees, reqListEmployees)

	if resListEmployees.Code != http.StatusOK {
		t.Errorf("ListEmployees() returned wrong status code: got %v want %v", resListEmployees.Code, http.StatusOK)
	}
}

func TestGetEmployee(t *testing.T) {
	store := NewEmployeeStore()

	employee := Employee{
		ID:       1,
		Name:     "John Doe",
		Position: "Software Engineer",
		Salary:   50000,
	}

	store.Employees[employee.ID] = employee

	router := mux.NewRouter()
	router.HandleFunc("/employee/{id}", store.GetEmployee)

	reqGetEmployee, _ := http.NewRequest("GET", "/employee/1", nil)
	resGetEmployee := httptest.NewRecorder()

	router.ServeHTTP(resGetEmployee, reqGetEmployee)

	if resGetEmployee.Code != http.StatusOK {
		t.Errorf("GetEmployee() returned wrong status code: got %v want %v", resGetEmployee.Code, http.StatusOK)
	}
}

func TestUpdateEmployee(t *testing.T) {
	store := NewEmployeeStore()

	employee := Employee{
		ID:       1,
		Name:     "John Doe",
		Position: "Software Engineer",
		Salary:   50000,
	}

	store.Employees[employee.ID] = employee

	updatedEmployee := Employee{
		ID:       1,
		Name:     "John Doe",
		Position: "Senior Software Engineer",
		Salary:   70000,
	}

	router := mux.NewRouter()
	router.HandleFunc("/employee/{id}", store.GetEmployee)

	updateEmployeeJSON, _ := json.Marshal(updatedEmployee)
	reqUpdateEmployee, _ := http.NewRequest("PUT", "/employee/1", bytes.NewBuffer(updateEmployeeJSON))
	reqUpdateEmployee.Header.Set("Content-Type", "application/json")
	resUpdateEmployee := httptest.NewRecorder()

	router.ServeHTTP(resUpdateEmployee, reqUpdateEmployee)

	if resUpdateEmployee.Code != http.StatusOK {
		t.Errorf("UpdateEmployee() returned wrong status code: got %v want %v", resUpdateEmployee.Code, http.StatusOK)
	}
}

func TestDeleteEmployee(t *testing.T) {
	store := NewEmployeeStore()

	employee := Employee{
		ID:       1,
		Name:     "John Doe",
		Position: "Software Engineer",
		Salary:   50000,
	}

	store.Employees[employee.ID] = employee

	router := mux.NewRouter()
	router.HandleFunc("/employee/{id}", store.GetEmployee)

	reqDeleteEmployee, _ := http.NewRequest("DELETE", "/employee/1", nil)
	resDeleteEmployee := httptest.NewRecorder()

	router.ServeHTTP(resDeleteEmployee, reqDeleteEmployee)
	if resDeleteEmployee.Code != http.StatusOK {
		t.Errorf("DeleteEmployee() returned wrong status code: got %v want %v", resDeleteEmployee.Code, http.StatusOK)
	}
}
