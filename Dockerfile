

FROM golang:alpine AS base

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY .. .


RUN apk add build-base

#COPY go.mod ./go.mod
#COPY go.sum ./go.sum

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest


#COPY . .

FROM base AS build

RUN swag init -g cmd/api/main.go output -o docs

RUN GOOS=linux go build -tags musl -ldflags "-w -s" -o bookmark_service cmd/api/main.go
RUN GOOS=linux go build -tags musl -ldflags "-w -s" -o bookmark_service_migrate cmd/migrate/main.go


FROM base AS test-exec

ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDE

RUN mkdir -p ${_outputdir} && \
    go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
    grep -v -E "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html

FROM scratch AS test
ARG _outputdir="/tmp/coverage"
COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

FROM alpine AS final
ARG app_name=app
ENV TZ=ASIA/Ho_Chi_Minh

WORKDIR /app

#COPY --from=build /opt/app/bookmark_service ./bookmark_service
#COPY --from=build /opt/app/docs ./docs
COPY --from=build /opt/app/bookmark_service /app/bookmark_service
COPY --from=build /opt/app/bookmark_service_migrate /app/bookmark_service_migrate
COPY --from=build /opt/app/docs /app/docs
COPY --from=build /opt/app/migrations /app/migrations

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD ["./bookmark_service"]
