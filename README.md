# hypertube-api
### .env
.env[developement | production].local
```
HYPERTUBE_ENV=development or production
API_PORT="3000"
POSTGRES_DB="hypertube"
POSTGRES_USER="docker"
POSTGRES_PASSWORD="secret"
POSTGRES_HOST="localhost"
POSTGRES_PORT="5432"
AUTH_42_CLIENT_ID="your 42 auth client id"
AUTH_42_SECRET="you 42 client secret"
AUTH_42_REDIRECT_URI="http://localhost:3000/api/auth/fortytwo/callback"
AUTH_42_AUTH_URI="https://api.intra.42.fr/oauth/authorize"
JWT_SECRET="your jwt secret"
MEDIA_HOST="localhost"
MEDIA_PORT="8080"
STATIC_FILES_PATH="./static for development or /static for production"
DOWNLOAD_FILES_PATH="./download for development or /download for production"
AES_KEY="your aes key"
JACKETT_API_KEY="your jackett api key"
JACKETT_HOST="localhost"
JACKETT_PORT="9117"
JACKETT_API_VERSION="v2.0"
```
