FROM golang:latest

COPY ./bin/app_linux_amd64 /usr/bin/weatherchecker-go

ENTRYPOINT [ "/usr/bin/weatherchecker-go" ]
