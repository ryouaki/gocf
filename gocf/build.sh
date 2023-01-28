./init.sh

go install

go env -w CGO_ENABLED=1

mkdir -p ./bin
mkdir -p ./bin/scripts

cp -rf ./scripts/* ./bin/scripts/*

go build -ldflags "-s -w" -o ./bin/gocf ./cmd/main.go 