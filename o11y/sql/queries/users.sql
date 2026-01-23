-- name: CreateUser :one
INSERT INTO users (id, name, age)
VALUES ($1, $2, $3)
RETURNING id, name, age; 
