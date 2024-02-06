INSTALL_CMD = sudo apt install sox
ifeq ($(OS),Windows_NT)
	INSTALL_CMD = winget install SoX
	INSTALL_CMD = scoop install sox
else
  UNAME_OS := $(shell uname -s)
  ifeq ($(UNAME_OS),Darwin)
		INSTALL_CMD = brew install sox
	endif
endif

build:
	go build .

run:
	go run .

depends:
	$(INSTALL_CMD)
	go mod tidy
