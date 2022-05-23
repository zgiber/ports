FROM golang:1.18 as builder
WORKDIR /ports
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
WORKDIR /ports/cmd
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -v -o ports

FROM alpine
COPY --from=builder /ports/cmd .
ENTRYPOINT [ "/ports" ]
CMD [ "-db_path","/db" ]
