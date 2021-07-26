BIN_API = api
BIN_MEDIA = media
BIN_TRANSCODER = transcoder
BIN_REGULATOR = regulator
BIN_HYPERDB = hyperdb
BIN_DIR = bin

all: $(BIN_API) $(BIN_MEDIA) $(BIN_TRANSCODER) $(BIN_REGULATOR)

start:
	docker-compose -f docker/compose/docker-compose.yml up -d psql
	sleep 3s
	pm2 start \
	$(BIN_DIR)/$(BIN_API) \
	$(BIN_DIR)/$(BIN_MEDIA) \
	$(BIN_DIR)/$(BIN_TRANSCODER) \

stop:
	pm2 stop \
	$(BIN_API) \
	$(BIN_MEDIA) \
	$(BIN_TRANSCODER) \

restart: stop start

test:
	go test -v ./pkg/...

flushdb:
	docker rm -f -v $(shell docker ps | grep 'psql' | awk -F ' ' '{print $$1}')
flushlogs:
	pm2 flush

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(BIN_API): | $(BIN_DIR)
	go build -v -o $(BIN_API) cmd/$(BIN_API)/main.go
	mv $(BIN_API) $(BIN_DIR)

$(BIN_MEDIA): | $(BIN_DIR)
	go build -v -o $(BIN_MEDIA) cmd/$(BIN_MEDIA)/main.go
	mv $(BIN_MEDIA) $(BIN_DIR)

$(BIN_TRANSCODER): | $(BIN_DIR)
	go build -v -o $(BIN_TRANSCODER) cmd/$(BIN_TRANSCODER)/main.go
	mv $(BIN_TRANSCODER) $(BIN_DIR)

$(BIN_REGULATOR): | $($(BIN_DIR))
	go build -v -o $(BIN_REGULATOR) cmd/$(BIN_REGULATOR)/main.go
	mv $(BIN_REGULATOR) $(BIN_DIR)

$(BIN_HYPERDB): | $($(BIN_DIR))
	go build -v -o $(BIN_HYPERDB) cmd/$(BIN_HYPERDB)/main.go
	mv $(BIN_HYPERDB) $(BIN_DIR)

clean:
	rm -rf $(BIN_DIR)/$(BIN_API) $(BIN_DIR)/$(BIN_MEDIA)
