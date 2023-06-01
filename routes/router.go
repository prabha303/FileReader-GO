package routes

import (
	"WORK/users/model"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

const FILE_NAME string = "employee.json"

// RouterConfig will provides the apis for http server
func RouterConfig() (router *httprouter.Router) {

	router = httprouter.New()

	router.GET("/v1/employees", GetEmployeesList)
	router.POST("/v1/employees", CreateEmployee)
	router.GET("/v1/employees/search/id/:id", FindByID)
	router.GET("/v1/employees/search/name/:name", FindByName)
	router.PUT("/v1/employee/update/id/:id", Update)

	return
}

func GetEmployeesList(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	file, err := os.Open(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	employees := []model.Employee{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		employee := model.Employee{}
		json.Unmarshal([]byte(scanner.Text()), &employee)
		employees = append(employees, employee)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	renderJSON(w, http.StatusOK, employees)
}

func FindByID(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {

	idStr := p.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, "INPUT ERROR ")
		return
	}

	file, err := os.Open(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	employee := model.Employee{}
	isExists := false
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		json.Unmarshal([]byte(scanner.Text()), &employee)

		if employee.ID == id {
			isExists = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if isExists {
		fmt.Println("Found Employee - ", employee)
		renderJSON(w, http.StatusOK, employee)
	} else {
		renderJSON(w, http.StatusBadRequest, "Employee Not Found")
	}
}

func FindByName(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {

	empName := p.ByName("name")
	if empName == "" {
		renderJSON(w, http.StatusBadRequest, "INPUT ERROR ")
		return
	}

	isExists := false

	file, err := os.Open(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	employee := model.Employee{}

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		json.Unmarshal([]byte(scanner.Text()), &employee)
		if strings.EqualFold(employee.Name, empName) {
			fmt.Println(employee.Name, empName)
			isExists = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if isExists {
		fmt.Println("Employee Found - ", employee)
		renderJSON(w, http.StatusOK, employee)
	} else {
		renderJSON(w, http.StatusBadRequest, "Employee Not Found")
	}
}

func CreateEmployee(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	fmt.Println("Init")

	employee := new(model.Employee)

	if !parseJSON(w, r.Body, &employee) {
		return
	}

	if employee == nil || employee.Name == "" {
		renderJSON(w, http.StatusBadRequest, "INPUT ERROR")
		return
	}

	// if !parseJSON(w, r.Body, &employee) {
	// 	return
	// }

	f, err := os.OpenFile(FILE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	newFsConfigBytes, _ := json.Marshal(employee)

	n, err := f.Write(newFsConfigBytes)
	if err != nil {
		fmt.Println(n, err)
	}

	if n, err = f.WriteString("\n"); err != nil {
		fmt.Println(n, err)
	}

	res := fmt.Sprintf("Employee Created Successfully : %s", employee.Name)

	renderJSON(w, http.StatusOK, res)
}

func Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	idStr := p.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, "INPUT ERROR")
		return
	}

	addressM := new(model.Address)

	if !parseJSON(w, r.Body, &addressM) {
		return
	}

	file, err := os.Open(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	employee := model.Employee{}
	isExists := false
	scanner := bufio.NewScanner(file)
	foundString := ""
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		foundString = scanner.Text()
		json.Unmarshal([]byte(foundString), &employee)
		if employee.ID == id {
			isExists = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if isExists {

		src, _ := ioutil.ReadFile(FILE_NAME)

		json.Unmarshal([]byte(foundString), &employee)
		employee.Address = addressM.Address

		empUpdate, _ := json.Marshal(employee)

		src = bytes.Replace(src, []byte(foundString), empUpdate, 1)
		if err = ioutil.WriteFile(FILE_NAME, src, 0666); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Updated Employee- ", employee)
		renderJSON(w, http.StatusOK, employee)
	} else {
		renderJSON(w, http.StatusBadRequest, "Employee Not Found")
	}
	//
}
