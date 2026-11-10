FROM node:24-alpine AS site-build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY index.html ./
COPY src ./src
RUN npm run build

FROM golang:1.27-alpine AS server-build
WORKDIR /app
COPY server/main.go ./main.go
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /site-server ./main.go

FROM traefik:v3.7.15
COPY --from=site-build /app/dist /srv/site
COPY --from=server-build /site-server /usr/local/bin/site-server
COPY traefik/ /etc/traefik/
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/site-server"]
