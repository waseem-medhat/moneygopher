FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

ADD bin/otps /usr/bin/otps

CMD ["otps"]
