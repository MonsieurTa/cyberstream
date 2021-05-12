BIN_API = api
BIN_MEDIA = media
BIN_DIR = bin

all: $(BIN_API) $(BIN_MEDIA)

start: all
	docker-compose up -d psql
	pm2 start bin/api bin/media

stop:
	pm2 stop api media

restart: stop start

flush:
	pm2 flush

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(BIN_API): | $(BIN_DIR)
	go build -o $(BIN_API) cmd/api/main.go
	mv $(BIN_API) $(BIN_DIR)
$(BIN_MEDIA): | $(BIN_DIR)
	go build -o $(BIN_MEDIA) cmd/media/main.go
	mv $(BIN_MEDIA) $(BIN_DIR)

clean:
	rm -rf $(BIN_DIR)/$(BIN_API) $(BIN_DIR)/$(BIN_MEDIA)
