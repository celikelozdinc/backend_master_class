-- name: CreateStudent :one
INSERT INTO students (name,number,nation) 
VALUES ($1, $2, $3) RETURNING *;

-- name: GetStudent :one
SELECT * FROM students
WHERE id = $1 LIMIT 1;

-- name: GetStudentForUpdate :one
SELECT * FROM students
WHERE id = $1 LIMIT 1
FOR UPDATE;

-- name: CreateGrade :one
INSERT INTO grades (student_id,grade) 
VALUES ($1, $2) RETURNING *;

-- name: GetGrade :one
SELECT * FROM grades
WHERE id = $1 LIMIT 1;


-- name: GetGradeForUpdate :one
SELECT * FROM grades
WHERE id = $1 LIMIT 1
FOR UPDATE;



-- name: GetGradeByStudentID :one
SELECT * FROM grades
WHERE student_id = $1 LIMIT 1;