package config

const (
	PORT       = ":3000"
	SQLITE_URI = "sqlite3"
	SQLITE_DB  = "test.db"

	POSTGRES_DB       = "hypertube"
	POSTGRES_USER     = "docker"
	POSTGRES_PASSWORD = "secret"
	POSTGRES_HOST     = "psql"

	PROVIDER_42_CLIENT_ID    = "f27a14bd15ac22eb8c548f3cb64147b40dfb85e17f071801b3fc2986ac4fe5e5"
	PROVIDER_42_SECRET       = "f98041bd1755f83a7a8b18892623b2ff8f591ae3f40de12e6c1ad1d24b5e665e"
	PROVIDER_42_REDIRECT_URI = "http://localhost:3000/api/auth/fortytwo/callback"
	PROVIDER_42_AUTH_URI     = "https://api.intra.42.fr/oauth/authorize"
)
