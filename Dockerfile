FROM golang:alpine AS build

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN go build -o bookmark-management cmd/api/main.go


#easy way
#CMD ["go","run","/opt/app/cmd/api/main.go"]
#CMD ["./bookmark-management"]

# difficult way
FROM alpine

WORKDIR /app

COPY --from=build /opt/app/bookmark-management /app/bookmark-management
COPY --from=build /opt/app/docs /app/docs

CMD ["/app/bookmark-management"]