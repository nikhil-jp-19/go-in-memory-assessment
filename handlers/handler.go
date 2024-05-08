package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Employee struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

type EmployeeStore struct {
	Employees map[int]Employee
	MU        sync.RWMutex
}

func NewEmployeeStore() *EmployeeStore {
	return &EmployeeStore{
		Employees: make(map[int]Employee),
	}
}

func (es *EmployeeStore) ListEmployees(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	es.MU.RLock()
	defer es.MU.RUnlock()

	employees := make([]Employee, 0)
	for id := 1; id <= len(es.Employees); id++ {
		employee, ok := es.Employees[id]
		if !ok {
			continue
		}
		if offset > 0 {
			offset--
			continue
		}
		employees = append(employees, employee)
		if len(employees) >= limit {
			break
		}
	}

	sort.Slice(employees, func(i, j int) bool {
		return employees[i].ID < employees[j].ID
	})

	json.NewEncoder(w).Encode(employees)
}

func (es *EmployeeStore) GetEmployee(w http.ResponseWriter, r *http.Request) {
	es.MU.RLock()
	defer es.MU.RUnlock()

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	employee, ok := es.Employees[id]
	if !ok {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(employee)
}

func (es *EmployeeStore) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	es.MU.Lock()
	defer es.MU.Unlock()

	id := len(es.Employees) + 1

	var employee Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	employee.ID = id

	es.Employees[id] = employee
	json.NewEncoder(w).Encode(employee)
}

func (es *EmployeeStore) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	es.MU.Lock()
	defer es.MU.Unlock()

	vars := mux.Vars(r)

	var updatedEmployee Employee
	if err := json.NewDecoder(r.Body).Decode(&updatedEmployee); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	if _, ok := es.Employees[id]; !ok {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	es.Employees[id] = updatedEmployee
	json.NewEncoder(w).Encode(updatedEmployee)
}

func (es *EmployeeStore) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	es.MU.Lock()
	defer es.MU.Unlock()

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	if _, ok := es.Employees[id]; !ok {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}
	delete(es.Employees, id)
	fmt.Fprintf(w, "Employee with ID %d deleted", id)
}
