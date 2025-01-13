CREATE DATABASE IF NOT EXISTS student_db;

USE student_db;

DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS branches;
DROP TABLE IF EXISTS student_names;

CREATE TABLE branches (
    branch_id INT AUTO_INCREMENT PRIMARY KEY,
    branch_name VARCHAR(50) NOT NULL,
    hod VARCHAR(50) NOT NULL
);

CREATE TABLE students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    student_id VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    branch_id INT,
    dob DATE NOT NULL,
    FOREIGN KEY (branch_id) REFERENCES branches(branch_id)
);

CREATE TABLE student_names (
    id INT PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL
);


