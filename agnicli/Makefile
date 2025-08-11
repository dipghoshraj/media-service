APP_NAME := agnicli
BUILD_DIR := bin

OS := $(shell uname -s)
EXT :=

ifeq ($(OS),Windows_NT)
    EXT := .exe
    INSTALL_PATH := $(USERPROFILE)\bin
    COPY := copy
else
    EXT :=
    INSTALL_PATH := /usr/local/bin
    COPY := cp
endif


BIN_PATH := $(BUILD_DIR)/$(APP_NAME)$(EXT)

.PHONY: all build install run clean

all: build

build:
	go build -o $(BIN_PATH) .

install: build
ifeq ($(OS),Windows_NT)
	$(COPY) $(BIN_PATH) $(INSTALL_PATH)\$(APP_NAME)$(EXT)
else
	sudo $(COPY) $(BIN_PATH) $(INSTALL_PATH)/$(APP_NAME)
endif

run: build
	$(BIN_PATH)

clean:
	rm -rf $(BUILD_DIR)
