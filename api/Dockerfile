# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.16.6-alpine AS build

WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod ./
COPY go.sum ./
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /server

### UPX
RUN apk add upx
RUN upx /server

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /

# Copy files that are required but not included in build
COPY --from=build /server /server
COPY --from=build /app/public /public
COPY --from=build /app/dev.env /
COPY --from=build /app/pro.env /

EXPOSE 8080

CMD [ "/server" ]