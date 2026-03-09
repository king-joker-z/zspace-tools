# ==================== 前端构建 ====================
FROM node:20-alpine AS frontend
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install --include=dev
COPY web/ .
RUN npx vite build

# ==================== 后端构建 ====================
FROM golang:1.22-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o zspace-tools .

# ==================== 最终镜像 ====================
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=backend /app/zspace-tools /usr/local/bin/
EXPOSE 8080
VOLUME ["/config", "/watch"]
ENV TZ=Asia/Shanghai
CMD ["zspace-tools"]
