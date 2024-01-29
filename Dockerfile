FROM golang:1.21.3 AS builder

## We create an /app directory in which
## we'll put all of our project code
RUN mkdir /app
ADD . /app
WORKDIR /app
## We want to build our application's binary executable
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./...

## the lightweight scratch image we'll
## run our application within
FROM alpine:latest AS production
## We have to copy the output from our
## builder stage to our production stage
ENV CONFIG_PHP_DBHOST host.docker.internal
ENV CONFIG_PHP_DBPORT 5432
ENV CONFIG_PHP_DBUSER busrai
ENV CONFIG_PHP_DBPASS 123
ENV CONFIG_PHP_DBNAME laravel_last
ENV CONFIG_PHP_DBTYPE pgsql
ENV CONFIG_PHP_DBSCHEME public
COPY --from=builder /app .
## we can then kick off our newly compiled
## binary exectuable!!
CMD ["./main"]