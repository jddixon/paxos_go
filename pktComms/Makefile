.PHONY = all

goFiles = p.pb.go client.pb.go paxos.pb.go 

%.pb.go: %.proto
	protoc --go_out=. $<
	@# perhaps superfluous
	@gofmt -w $@

all: $(goFiles)
