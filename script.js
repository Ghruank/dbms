document.addEventListener("DOMContentLoaded", function() {
    fetchBranches();
    fetchStudents();

    document.getElementById("registrationForm").addEventListener("submit", function(event) {
        event.preventDefault();
        const id = document.getElementById("id").value;
        if (id) {
            updateStudent(id);
        } else {
            registerStudent();
        }
    });

    document.getElementById("fetchForm").addEventListener("submit", function(event) {
        event.preventDefault();
        const studentId = document.getElementById("fetchStudentId").value;
        fetchStudentById(studentId);
    });

    document.getElementById("branchForm").addEventListener("submit", function(event) {
        event.preventDefault();
        addBranch();
    });

    document.getElementById("branch").addEventListener("change", function() {
        const selectedBranch = this.value;
        const branches = JSON.parse(localStorage.getItem("branches"));
        const branch = branches.find(branch => branch.branch_name === selectedBranch);
        document.getElementById("hod").value = branch ? branch.hod : "";
    });
});

function fetchBranches() {
    fetch("http://localhost:8080/branches")
        .then(response => response.json())
        .then(data => {
            const branchSelect = document.getElementById("branch");
            branchSelect.innerHTML = '<option value="" disabled selected>Branch</option>'; // Clear existing options and add default option
            localStorage.setItem("branches", JSON.stringify(data));
            data.forEach(branch => {
                const option = document.createElement("option");
                option.value = branch.branch_name;
                option.textContent = branch.branch_name;
                branchSelect.appendChild(option);
            });
        })
        .catch(error => {
            console.error("Error fetching branches:", error);
        });
}

function fetchStudents() {
    fetch("http://localhost:8080/students")
        .then(response => response.json())
        .then(data => {
            const tableBody = document.getElementById("studentsTable").getElementsByTagName("tbody")[0];
            tableBody.innerHTML = "";
            data.forEach(student => {
                const row = tableBody.insertRow();
                row.insertCell(0).textContent = student.student_id;
                row.insertCell(1).textContent = student.name;
                row.insertCell(2).textContent = student.branch;
                row.insertCell(3).textContent = student.hod;
                row.insertCell(4).textContent = student.dob;
                row.insertCell(5).textContent = student.age;
                const actionsCell = row.insertCell(6);
                const editButton = document.createElement("button");
                editButton.textContent = "Edit";
                editButton.onclick = function() {
                    fetchStudentById(student.student_id);
                };
                actionsCell.appendChild(editButton);
            });
        })
        .catch(error => {
            console.error("Error fetching students:", error);
        });
}

function fetchStudentById(studentId) {
    fetch(`http://localhost:8080/student?id=${studentId}`)
        .then(response => response.json())
        .then(student => {
            document.getElementById("id").value = student.student_id;
            document.getElementById("name").value = student.name;
            document.getElementById("branch").value = student.branch;
            document.getElementById("dob").value = student.dob;
            document.getElementById("submitButton").textContent = "Update";
        })
        .catch(error => {
            console.error("Error fetching student:", error);
        });
}

function registerStudent() {
    const form = document.getElementById("registrationForm");
    const formData = new FormData(form);

    fetch("http://localhost:8080/register", {
        method: "POST",
        body: formData
    })
    .then(response => response.text())
    .then(message => {
        alert(message);
        fetchStudents();
        form.reset(); // Clear the form inputs
        document.getElementById("submitButton").textContent = "Register";
    })
    .catch(error => {
        console.error("Error:", error);
    });
}

function updateStudent(studentId) {
    const form = document.getElementById("registrationForm");
    const formData = new FormData(form);

    fetch(`http://localhost:8080/update`, {
        method: "POST",
        body: formData
    })
    .then(response => response.text())
    .then(message => {
        alert(message);
        fetchStudents();
        form.reset(); // Clear the form inputs
        document.getElementById("submitButton").textContent = "Register";
    })
    .catch(error => {
        console.error("Error:", error);
    });
}

function addBranch() {
    const form = document.getElementById("branchForm");
    const formData = new FormData(form);

    fetch("http://localhost:8080/addBranch", {
        method: "POST",
        body: formData
    })
    .then(response => response.text())
    .then(message => {
        alert(message);
        fetchBranches();
        form.reset(); // Clear the form inputs
    })
    .catch(error => {
        console.error("Error:", error);
    });
}
