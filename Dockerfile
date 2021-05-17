## Dockerfile for the Sea of Thieves Bot
## 
FROM        archlinux
LABEL       maintainer="wn@neessen.net"
RUN         pacman -Syu --noconfirm --noprogressbar
RUN         /usr/bin/groupadd -r sotbot && /usr/bin/useradd -r -g sotbot -c "Sea of Thieves Bot" -m -s /bin/bash -d /opt/sotbot sotbot
COPY        ["LICENSE", "README.md", "/opt/sotbot/"]
COPY        ["bin", "/opt/sotbot/bin"]
COPY        ["media", "/opt/sotbot/media"]
COPY        ["config", "/opt/sotbot/config"]
RUN         chown -R sotbot:sotbot /opt/sotbot
WORKDIR     /opt/sotbot
USER        sotbot
VOLUME      ["/opt/sotbot/config"]
ENTRYPOINT  ["/opt/sotbot/bin/sotbot", "-c", "/opt/sotbot/config"]
