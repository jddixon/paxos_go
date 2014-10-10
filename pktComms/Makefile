p.pb.go: p.proto client.proto paxos.proto 
	mkdir -p _pb
	protoc --go_out=_pb $^
	gofmt -w _pb
	mv _pb/* .
	rm -rf _pb
