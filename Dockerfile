FROM golang:latest

COPY ./bin/app_linux_amd64 /usr/bin/weatherchecker-go

CMD [ "/usr/bin/weatherchecker-go" ]
