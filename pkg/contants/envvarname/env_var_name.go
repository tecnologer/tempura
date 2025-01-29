package envvarname

const (
	DBHost     = "TEMPURA_DB_HOST"
	DBPort     = "TEMPURA_DB_PORT"
	DBUsername = "TEMPURA_DB_USERNAME"
	DBPassword = "TEMPURA_DB_PASSWORD" //nolint:gosec // This is not a real password
	DBName     = "TEMPURA_DB_NAME"
	DBSSLMode  = "TEMPURA_DB_SSL_MODE"
)
