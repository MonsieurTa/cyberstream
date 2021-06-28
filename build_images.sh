docker build -t mrta/hypertube-api:latest -f cmd/api/Dockerfile .
docker build -t mrta/hypertube-media:latest -f cmd/media/Dockerfile .
docker build -t mrta/hypertube-transcoder:latest -f cmd/transcoder/Dockerfile .

docker push mrta/hypertube-api:latest
docker push mrta/hypertube-media:latest
docker push mrta/hypertube-transcoder:latest
