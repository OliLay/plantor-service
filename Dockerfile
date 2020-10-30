FROM golang as builder

CMD mkdir /plantor-service
WORKDIR /plantor-service

COPY go.mod .
COPY go.sum .
COPY /influx /plantor-service/influx
COPY /mqtt /plantor-service/mqtt
COPY plantor.go /plantor-service/plantor.go

RUN CGO_ENABLED=0 go build

FROM alpine
COPY --from=builder /plantor-service /plantor-service

ENTRYPOINT ["/plantor-service/plantor"]