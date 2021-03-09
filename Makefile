all: rmgem rmgem.arm

rmgem: $(wildcard *.go)
	go build -o rmgem

rmgem.arm: $(wildcard *.go)
	GOOS=linux GOARCH=arm go build -o rmgem.arm

.PHONY: clean
clean:
	rm rmgem rmgem.arm
