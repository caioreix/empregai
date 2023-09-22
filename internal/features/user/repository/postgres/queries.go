package postgres

const (
	createUserQuery = `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING *
	`

	updateUserQuery = `
		UPDATE users
		SET name = COALESCE(NULLIF($2, ''), name),
			email = COALESCE(NULLIF($3, ''), email),
			password = COALESCE(NULLIF($4, ''), password)
		WHERE id = $1
		RETURNING *
	`

	deleteUserQuery = `DELETE FROM users WHERE id = $1`

	getUserByIDQuery = `
		SELECT id, name, email, password, created_at, updated_at, last_login
		FROM users
		WHERE id = $1
	`

	getUserByEmailQuery = `
		SELECT id, name, email, password, created_at, updated_at, last_login
		FROM users
		WHERE email = $1
	`

	getUsersCountQuery = `SELECT COUNT(id) FROM users`

	getAllUsersQuery = `
		SELECT id, name, email, password, created_at, updated_at, last_login
		FROM users
		ORDER BY COALESCE(NULLIF($1, ''), name)
		OFFSET $2
		LIMIT $3
	`
)
