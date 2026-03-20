#! /bin/sh

rm -r build/*

CGO_ENABLED=0 go build -o build/easyjet cmd/server/main.go

npm install -C frontend
npm run build -C frontend

mv frontend/dist build/web