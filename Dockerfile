FROM golang:1.18-buster AS builder
WORKDIR /app
COPY go.* ./
COPY settings.json ./
RUN go mod download
COPY *.go ./
RUN go build -o /json-excel-conv
# Create a new release build stage
FROM gcr.io/distroless/base-debian10
# Set the working directory to the root directory path
WORKDIR /
# Copy over the binary built from the previous stage
COPY --from=builder /json-excel-conv /json-excel-conv
COPY --from=builder /app/settings.json /settings.json
EXPOSE 8080
ENTRYPOINT ["/json-excel-conv"]
