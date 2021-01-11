GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GODIR=n_puzzle
TARGET=n-puzzle

all: build

build:
	cd $(GODIR) && $(GOBUILD) -o $(TARGET)
	mv ./$(GODIR)/$(TARGET) ./

clean:
	$(GOCLEAN)
	rm $(TARGET)

re: clean build