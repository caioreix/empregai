package postgres

const (
	createUserQuery = `
		INSERT INTO users (email, password, role)
		VALUES ($1, $2, $3)
		RETURNING *
	`

	updateUserQuery = `
		UPDATE users
		SET email = COALESCE(NULLIF($2, ''), email),
			password = COALESCE(NULLIF($3, ''), password),
			role = COALESCE(NULLIF($4, ''), role)
		WHERE id = $1
		RETURNING *
	`

	deleteUserQuery = `DELETE FROM users WHERE id = $1`

	getUserByIDQuery = `
		SELECT id, email, password, role, created_at, updated_at, last_login
		FROM users
		WHERE id = $1
	`

	getUserByEmailQuery = `
		SELECT id, email, password, role, created_at, updated_at, last_login
		FROM users
		WHERE email = $1
	`

	getUsersCountQuery = `SELECT COUNT(id) FROM users`

	getAllUsersQuery = `
		SELECT id, email, password, role, created_at, updated_at, last_login
		FROM users
		ORDER BY COALESCE(NULLIF($1, ''), email)
		OFFSET $2
		LIMIT $3
	`
)
