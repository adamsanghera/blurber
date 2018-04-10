rm ./blurber.pb.go > /dev/null 2>&1
protoc ./blurber.proto --go_out=plugins=grpc:./