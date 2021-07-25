# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.16.6-alpine AS build

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /server

### Stage 2
RUN apk add upx
RUN upx /server
FROM scratch
WORKDIR /
COPY --from=build /server /server
COPY --from=build /app/public /public
COPY --from=build /app/dev.env /
COPY --from=build /app/pro.env /

EXPOSE 8080

CMD [ "/server" ]