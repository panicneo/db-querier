FROM golang:alpine AS build-stage-backend
WORKDIR /go/src/db-querier
COPY ./backend .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags=jsoniter -o server .

FROM node:alpine AS build-stage-frontend
WORKDIR /db-querier
COPY ./frontend .
RUN npm install --registry=https://registry.npm.taobao.org
RUN npm run build

FROM alpine
ENV TZ=Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
  && apk --no-cache --update add tzdata ca-certificates

WORKDIR /server
COPY --from=build-stage-backend /go/src/db-querier/server .
COPY --from=build-stage-backend /go/src/db-querier/configs ./configs
COPY --from=build-stage-frontend /db-querier/dist ./dist

CMD ["./server"]
