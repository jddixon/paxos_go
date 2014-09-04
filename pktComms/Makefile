p.pb.go: p.proto
	mkdir -p _pb
	protoc --go_out=_pb $<
	cat _pb/$@\
	|sed s/Newmsg/newMsg/g\
	|gofmt >$@
	rm -rf _pb
