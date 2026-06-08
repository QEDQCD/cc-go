.PHONY: all build web clean run dev test build-mac-app release

APP_NAME = cc-go
WEB_DIR = web
DIST_DIR = dist
LDFLAGS = -s -w
GOPROXY ?= https://goproxy.cn,direct
export GOPROXY

all: web build

web:
	cd $(WEB_DIR) && npm install && npm run build

build:
	go build -ldflags "-H windowsgui -s -w" -o $(APP_NAME) ./cmd/cc-go/

build-linux:
	cd $(WEB_DIR) && npm run build
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(APP_NAME)-linux ./cmd/cc-go/

build-mac:
	cd $(WEB_DIR) && npm run build
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $(APP_NAME)-mac ./cmd/cc-go/

build-mac-app:
	cd $(WEB_DIR) && npm run build
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o $(APP_NAME)-darwin-arm64 ./cmd/cc-go/
	@mkdir -p "$(APP_NAME).app/Contents/MacOS"
	@mkdir -p "$(APP_NAME).app/Contents/Resources"
	@cp $(APP_NAME)-darwin-arm64 "$(APP_NAME).app/Contents/MacOS/cc-go"
	@cp cmd/cc-go/Info.plist "$(APP_NAME).app/Contents/"
	@cp cmd/cc-go/cc-go.icns "$(APP_NAME).app/Contents/Resources/"

build-win:
	cd $(WEB_DIR) && npm run build
	GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui -s -w" -o $(APP_NAME).exe ./cmd/cc-go/

run:
	go run -ldflags "-H windowsgui -s -w" ./cmd/cc-go/

dev:
	cd $(WEB_DIR) && npm run dev &
	go run -ldflags "-H windowsgui -s -w" ./cmd/cc-go/

clean:
	rm -rf cmd/cc-go/web-dist $(APP_NAME) $(APP_NAME)-linux $(APP_NAME)-mac $(APP_NAME).exe $(APP_NAME).app $(APP_NAME)-darwin-*

test:
	go test ./... -v -count=1 -timeout 60s

release: web
	@mkdir -p $(DIST_DIR)
	cd $(WEB_DIR) && npm run build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 ./cmd/cc-go/
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui $(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe ./cmd/cc-go/
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 ./cmd/cc-go/
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 ./cmd/cc-go/
	@cd $(DIST_DIR) && sha256sum $(APP_NAME)-* > SHA256SUMS
	@echo "Release binaries built in $(DIST_DIR)/"