FROM golang:1.20 as builder

# STAGE: BUILD

WORKDIR /app

# latest commit hash/title is used for building the project
COPY .git .git

COPY go.mod .
COPY go.sum .
COPY Makefile .

COPY assets assets
COPY internal internal
COPY cmd cmd

RUN make dist

# STAGE: TARGET

FROM alpine:latest
RUN apk add ffmpeg

RUN addgroup -S ary && adduser -S ary -G ary

USER ary

WORKDIR /app
COPY --from=builder /app/assets /app/assets
COPY --from=builder /app/aryzona /app/aryzona

ENTRYPOINT ["/app/aryzona"]
