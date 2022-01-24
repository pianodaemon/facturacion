# Building stage
FROM golang:1.15-buster as builder

LABEL MAINTAINER="Edwin Plauchu <pianodaemon@gmail.com>"

WORKDIR /go/src/immortalcrab.com/interview

COPY Makefile go.mod ./

RUN go mod download

COPY pkg pkg
COPY internal internal
COPY cmd cmd

RUN make clean build

# Final image
FROM debian:bullseye
LABEL MAINTAINER="Edwin Plauchu <pianodaemon@gmail.com>"

ENV APP_DIR=/
COPY --from=builder /interview $APP_DIR
COPY scripts/run_service.sh $APP_DIR
WORKDIR $APP_DIR

EXPOSE 10100

CMD ["/run_service.sh"]
