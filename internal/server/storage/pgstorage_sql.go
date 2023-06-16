package storage

const (
	// Users.
	_sqlCreateTableUser = `
		CREATE TABLE IF NOT EXISTS users (
			id         serial                   NOT NULL,
			username   text                     NOT NULL,
			password   text                     NOT NULL,
			created_at timestamp with time zone NOT NULL DEFAULT now(),
			
			PRIMARY KEY (id),
			UNIQUE(username)
		);
		`
	_sqlCreateUser = `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
		RETURNING id
		`
	_sqlFindUser = `
		SELECT id, password FROM users
		WHERE username = $1
		`
)
