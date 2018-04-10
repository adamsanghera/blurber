# Remove the dead weeds, if they exist
rm -rf ./dist/*/* > /dev/null 2>&1

# Make our path
mkdir dist > /dev/null 2>&1
mkdir dist/blurb > /dev/null 2>&1
mkdir dist/user > /dev/null 2>&1
mkdir dist/subscription > /dev/null 2>&1

# Complie the files
protoc ./blurb.proto --go_out=plugins=grpc:./dist/blurb/
protoc ./user.proto --go_out=plugins=grpc:./dist/user/
protoc ./subscription.proto --go_out=plugins=grpc:./dist/subscription/