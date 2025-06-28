FROM golang:1.24.4-alpine3.22 AS builder

WORKDIR /src
COPY . .

RUN apk add make

RUN make install
RUN make build


FROM alpine:3.22 AS runner

COPY --from=builder /src/build/baitomebot /baitomebot

ENTRYPOINT [ "/baitomebot" ]
