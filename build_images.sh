docker build -t mrta/hypertube-api:latest -f cmd/api/Dockerfile .
docker build -t mrta/hypertube-media:latest -f cmd/media/Dockerfile .
docker build -t mrta/hypertube-transcoder:latest -f cmd/transcoder/Dockerfile .
