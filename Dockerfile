FROM alpine as server_source
WORKDIR /app
COPY . .
RUN rm -rf ./ui

FROM golang:alpine as build_server
RUN apk add npm make gcc g++; mkdir -p /app
WORKDIR /app
# Copy lockfiles first to make sure that
# download cache is only busted if necessary
COPY go.* ./
RUN go mod download
# Then copy everything else (but not ui, 
# so that changes in ui do not cause rebuild)
COPY --from=server_source /app .
RUN go build -o /app/bin/server /app/cmd/server/main.go


FROM node:alpine as build_ui
RUN mkdir -p /app/ui
WORKDIR /app/ui
# First, copy package.json and package-lock.json
# so that we keep node_modules unless one of these changes
COPY ./ui/package* ./
RUN npm install
# Then copy everything else and build app
COPY ./ui ./
RUN npm run build


FROM alpine as prod
RUN mkdir -p /app/ui /app/bin; apk add libstdc++
COPY --from=build_ui /app/ui/build/ /app/ui/build/
COPY --from=build_server /app/bin/server /app/bin/server
WORKDIR /app
CMD ./bin/server