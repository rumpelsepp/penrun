FROM debian:latest

RUN apt-get update && apt-get install -y bats jq zstd git shellcheck shfmt reuse make

WORKDIR /code
ENTRYPOINT ["bash"]
