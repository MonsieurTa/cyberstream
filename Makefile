BIN_API = api
BIN_MEDIA = media
BIN_DIR = bin

all: $(BIN_API) $(BIN_MEDIA)

$(BIN_DIR):
	mkdir -p bin

$(BIN_API): | $(BIN_DIR)
	go build -o $(BIN_API) cmd/api/main.go
	mv $(BIN_API) $(BIN_DIR)
$(BIN_MEDIA): | $(BIN_DIR)
	go build -o $(BIN_MEDIA) cmd/media/main.go
	mv $(BIN_MEDIA) $(BIN_DIR)

clean:
	rm -rf $(BIN_DIR)/$(BIN_API) $(BIN_DIR)/$(BIN_MEDIA)
