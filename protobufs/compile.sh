# Remove the dead weeds, if they exist
rm -rf ./dist > /dev/null 2>&1

# Make our paths
mkdir dist > /dev/null 2>&1
mkdir dist/blurb > /dev/null 2>&1
mkdir dist/user > /dev/null 2>&1
mkdir dist/subscription > /dev/null 2>&1
mkdir dist/common > /dev/null 2>&1

# Complie the files

protoc --go_out=plugins=grpc:$GOPATH/src \
       ./common.proto

protoc --go_out=plugins=grpc:$GOPATH/src/github.com/adamsanghera/blurber/protobufs/dist/blurb \
       ./blurb.proto

protoc --go_out=plugins=grpc:$GOPATH/src/github.com/adamsanghera/blurber/protobufs/dist/user \
       ./user.proto

protoc --go_out=plugins=grpc:$GOPATH/src/github.com/adamsanghera/blurber/protobufs/dist/subscription \
       ./subscription.proto