module github.com/MonsieurTa/hypertube

go 1.15

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/anacrolix/log v0.8.0
	github.com/anacrolix/torrent v1.26.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.2.0
	github.com/google/uuid v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.3.0
	github.com/stretchr/testify v1.7.0
	github.com/webtor-io/go-jackett v0.0.0-20201110160721-0d56a2f41070
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.21.6
)

replace github.com/webtor-io/go-jackett => github.com/MonsieurTa/go-jackett v0.0.4
