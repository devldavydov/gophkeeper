//nolint:gosec // OK
package storage

const (
	// Users.
	_sqlCreateTableUser = `
		CREATE TABLE IF NOT EXISTS users (
			id         bigserial                NOT NULL,
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
	// Secrets.
	_sqlCreateTableSecret = `
		CREATE TABLE IF NOT EXISTS secrets (
			user_id     bigint    NOT NULL,
			type        int       NOT NULL,
			name        text      NOT NULL,
			meta        text,
			version     bigint    NOT NULL,
			payload_raw bytea     NOT NULL,

			PRIMARY KEY (user_id, name),
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
		`
	_sqlCreateSecret = `
		INSERT INTO secrets (user_id, type, name, meta, version, payload_raw)
		VALUES ($1, $2, $3, $4, $5, $6);
		`
)
