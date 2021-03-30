module github.com/MonsieurTa/hypertube

go 1.15

require (
	db v1.0.0
	gorm.io/driver/sqlite v1.1.4 // indirect
	gorm.io/gorm v1.21.6 // indirect
	models v1.0.0 // indirect
	repo v1.0.0 // indirect
)

replace (
	db v1.0.0 => ./pkg/internal/db
	models v1.0.0 => ./pkg/internal/models
	repo v1.0.0 => ./pkg/internal/repo
)
