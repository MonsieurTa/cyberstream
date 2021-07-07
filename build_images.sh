DOCKER_BUILDKIT=1 docker build -t mrta/hypertube-api:latest -f cmd/api/Dockerfile .
DOCKER_BUILDKIT=1 docker build -t mrta/hypertube-media:latest -f cmd/media/Dockerfile .
DOCKER_BUILDKIT=1 docker build -t mrta/hypertube-transcoder:latest -f cmd/transcoder/Dockerfile .
