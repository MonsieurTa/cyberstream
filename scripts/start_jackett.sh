mkdir -p $HOME/Projects/jackett/config
mkdir -p $HOME/Projects/jackett/download

docker run -d \
    --name=jackett \
    -e PUID=1000 \
    -e PGID=1000 \
    -e TZ=Europe/London \
    -p 9117:9117 \
    -v $HOME/Projects/jackett/config:/config \
    -v $HOME/Projects/jackett/download:/downloads \
    --restart unless-stopped \
    ghcr.io/linuxserver/jackett
