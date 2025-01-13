package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "ghruank:ghruankonsql@tcp(127.0.0.1:3306)/student_db")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection established")
}

func calculateAge(dob string) int {
	layout := "2006-01-02"
	birthdate, _ := time.Parse(layout, dob)
	today := time.Now()
	age := today.Year() - birthdate.Year()
	if today.YearDay() < birthdate.YearDay() {
		age--
	}
	return age
}

func registerStudent(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to register student")
	if r.Method == "POST" {
		studentId := r.FormValue("studentId")
		name := r.FormValue("name")
		branch := r.FormValue("branch")
		dob := r.FormValue("dob")

		log.Printf("Student ID: %s, Name: %s, Branch: %s, DOB: %s\n", studentId, name, branch, dob)

		var branchId int
		err := db.QueryRow("SELECT branch_id FROM branches WHERE branch_name = ?", branch).Scan(&branchId)
		if err != nil {
			http.Error(w, "Branch not found", http.StatusNotFound)
			log.Println("Branch not found:", branch, "Error:", err)
			return
		}

		log.Printf("Inserting student: StudentID=%s, Name=%s, BranchID=%d, DOB=%s", studentId, name, branchId, dob)
		result, err := db.Exec("INSERT INTO students (student_id, name, branch_id, dob) VALUES (?, ?, ?, ?)", studentId, name, branchId, dob)
		if err != nil {
			http.Error(w, "Failed to register student", http.StatusInternalServerError)
			log.Println("Failed to register student:", err)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Failed to retrieve last insert ID", http.StatusInternalServerError)
			log.Println("Failed to retrieve last insert ID:", err)
			return
		}

		// Split the name into first name and last name
		nameParts := strings.SplitN(name, " ", 2)
		firstName := nameParts[0]
		lastName := ""
		if len(nameParts) > 1 {
			lastName = nameParts[1]
		}

		log.Printf("Inserting student name: FirstName=%s, LastName=%s", firstName, lastName)
		_, err = db.Exec("INSERT INTO student_names (id, first_name, last_name) VALUES (?, ?, ?)", id, firstName, lastName)
		if err != nil {
			http.Error(w, "Failed to insert student name", http.StatusInternalServerError)
			log.Println("Failed to insert student name:", err)
			return
		}

		age := calculateAge(dob)
		fmt.Fprintf(w, "Student registered successfully. Age: %d", age)
		log.Println("Student registered successfully")
	}
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get students")
	rows, err := db.Query("SELECT s.student_id, s.name, b.branch_name, b.hod, s.dob FROM students s JOIN branches b ON s.branch_id = b.branch_id")
	if err != nil {
		http.Error(w, "Failed to retrieve students", http.StatusInternalServerError)
		log.Println("Failed to retrieve students:", err)
		return
	}
	defer rows.Close()

	var students []map[string]interface{}
	for rows.Next() {
		var studentId, name, branch, hod, dob string
		if err := rows.Scan(&studentId, &name, &branch, &hod, &dob); err != nil {
			http.Error(w, "Failed to scan student", http.StatusInternalServerError)
			log.Println("Failed to scan student:", err)
			return
		}
		age := calculateAge(dob)
		student := map[string]interface{}{
			"student_id": studentId,
			"name":       name,
			"branch":     branch,
			"hod":        hod,
			"dob":        dob,
			"age":        age,
		}
		students = append(students, student)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func getStudentByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get student by ID")
	studentId := r.URL.Query().Get("id")
	if studentId == "" {
		http.Error(w, "Missing student ID", http.StatusBadRequest)
		return
	}

	var student struct {
		StudentID string `json:"student_id"`
		Name      string `json:"name"`
		Branch    string `json:"branch"`
		Dob       string `json:"dob"`
	}

	err := db.QueryRow("SELECT s.student_id, s.name, b.branch_name, s.dob FROM students s JOIN branches b ON s.branch_id = b.branch_id WHERE s.student_id = ?", studentId).Scan(&student.StudentID, &student.Name, &student.Branch, &student.Dob)
	if err != nil {
		http.Error(w, "Failed to retrieve student", http.StatusInternalServerError)
		log.Println("Failed to retrieve student:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to update student")
	if r.Method == "POST" {
		studentId := r.FormValue("studentId")
		name := r.FormValue("name")
		branch := r.FormValue("branch")
		dob := r.FormValue("dob")

		log.Printf("Student ID: %s, Name: %s, Branch: %s, DOB: %s\n", studentId, name, branch, dob)

		var branchId int
		err := db.QueryRow("SELECT branch_id FROM branches WHERE branch_name = ?", branch).Scan(&branchId)
		if err != nil {
			http.Error(w, "Branch not found", http.StatusNotFound)
			log.Println("Branch not found:", branch, "Error:", err)
			return
		}

		log.Printf("Updating student: StudentID=%s, Name=%s, BranchID=%d, DOB=%s", studentId, name, branchId, dob)
		_, err = db.Exec("UPDATE students SET name = ?, branch_id = ?, dob = ? WHERE student_id = ?", name, branchId, dob, studentId)
		if err != nil {
			http.Error(w, "Failed to update student", http.StatusInternalServerError)
			log.Println("Failed to update student:", err)
			return
		}

		// Split the name into first name and last name
		nameParts := strings.SplitN(name, " ", 2)
		firstName := nameParts[0]
		lastName := ""
		if len(nameParts) > 1 {
			lastName = nameParts[1]
		}

		log.Printf("Updating student name: FirstName=%s, LastName=%s", firstName, lastName)
		_, err = db.Exec("UPDATE student_names SET first_name = ?, last_name = ? WHERE id = (SELECT id FROM students WHERE student_id = ?)", firstName, lastName, studentId)
		if err != nil {
			http.Error(w, "Failed to update student name", http.StatusInternalServerError)
			log.Println("Failed to update student name:", err)
			return
		}

		fmt.Fprintf(w, "Student updated successfully")
		log.Println("Student updated successfully")
	}
}

func addBranch(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to add branch")
	if r.Method == "POST" {
		branchName := r.FormValue("branchName")
		hod := r.FormValue("hod")

		log.Printf("Branch Name: %s, HOD: %s\n", branchName, hod)

		_, err := db.Exec("INSERT INTO branches (branch_name, hod) VALUES (?, ?)", branchName, hod)
		if err != nil {
			http.Error(w, "Failed to add branch", http.StatusInternalServerError)
			log.Println("Failed to add branch:", err)
			return
		}
		fmt.Fprintf(w, "Branch added successfully")
		log.Println("Branch added successfully")
	}
}

func getBranches(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get branches")
	rows, err := db.Query("SELECT branch_name, hod FROM branches")
	if err != nil {
		http.Error(w, "Failed to retrieve branches", http.StatusInternalServerError)
		log.Println("Failed to retrieve branches:", err)
		return
	}
	defer rows.Close()

	var branches []map[string]interface{}
	for rows.Next() {
		var branchName, hod string
		if err := rows.Scan(&branchName, &hod); err != nil {
			http.Error(w, "Failed to scan branch", http.StatusInternalServerError)
			log.Println("Failed to scan branch:", err)
			return
		}
		branch := map[string]interface{}{
			"branch_name": branchName,
			"hod":         hod,
		}
		branches = append(branches, branch)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branches)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/register", registerStudent)
	http.HandleFunc("/students", getStudents)
	http.HandleFunc("/student", getStudentByID)
	http.HandleFunc("/update", updateStudent)
	http.HandleFunc("/addBranch", addBranch)
	http.HandleFunc("/branches", getBranches)
	http.Handle("/", http.FileServer(http.Dir("."))) // Serve static files from the current directory
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
