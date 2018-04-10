rm ./blurb.pb.go > /dev/null 2>&1
rm ./user.pb.go > /dev/null 2>&1
protoc ./blurb.proto --go_out=plugins=grpc:./
protoc ./user.proto --go_out=plugins=grpc:./