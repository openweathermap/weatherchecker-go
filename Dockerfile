FROM centurylink/ca-certs

COPY ./bin/app_linux_amd64 /usr/bin/app

ENTRYPOINT [ "/usr/bin/app" ]
