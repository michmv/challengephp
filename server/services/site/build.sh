#!/bin/bash

if [[ ! -d "/go/bin" ]]; then
    # если папка bin не существует устанавливаем в систему go
    cd /
    echo "Собираем golang" >> /app/out.log
    wget https://go.dev/dl/go1.24.5.linux-amd64.tar.gz
    tar -C /go -xzf go1.24.5.linux-amd64.tar.gz
    rm -f go1.24.5.linux-amd64.tar.gz

    CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

    echo "Golang собран" >> /app/out.log
fi

while true; do
    echo "Скрипт работает: $(date)"
    sleep 10
done
