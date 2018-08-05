FROM golang:1.10 AS builder

# Copy the code from the host and compile it
ADD ./ /go/src/github.com/kypseli/todo-api/
WORKDIR /go/src/github.com/kypseli/todo-api/

RUN go install github.com/kypseli/todo-api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM scratch
COPY --from=builder /app ./
ENTRYPOINT ["./app"]