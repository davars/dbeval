-- name: InsertAuthors :copyfrom
INSERT INTO authors (id, name)
VALUES ($1, $2);

-- name: InsertArticles :copyfrom
INSERT INTO articles (id, title, body, published_at)
VALUES ($1, $2, $3, $4);

-- name: FindAuthorByID :one
SELECT id, name
FROM authors
WHERE id = $1
LIMIT 1;

-- name: FindAuthorsByName :many
SELECT id, name
FROM authors
WHERE name = $1;

-- name: RecentArticles :many
SELECT id, title, body, published_at
FROM articles
ORDER BY published_at DESC
LIMIT $1;
