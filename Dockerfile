# STAGE: BUILD
FROM golang:1.23 AS builder

WORKDIR /app

# latest commit hash/title is used for building the project
COPY .git .git

COPY go.mod .
COPY go.sum .
COPY Makefile .

COPY assets assets
COPY migration migration
COPY internal internal
COPY cmd cmd

RUN go mod download

RUN make dist

# STAGE: TARGET

FROM alpine:3

RUN addgroup -S ary && adduser -S ary -G ary

USER ary

WORKDIR /app
COPY --from=builder /app/migration /app/migration
COPY --from=builder /app/assets /app/assets
COPY --from=builder /app/aryzona /app/aryzona

ENTRYPOINT ["/app/aryzona"]
