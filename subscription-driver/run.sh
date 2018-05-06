echo "Cleaning up old logs"

rm ./*.txt > /dev/null

echo "Building the binary"

go build ./

echo "Binary built, executing procs..."

SUB_PORT=6000 LEADER_ADDRESS=0 ./subscription-driver &> leader.txt &disown
SUB_PORT=6010 LEADER_ADDRESS=127.0.0.1:6001 ./subscription-driver &> follower1.txt &disown
SUB_PORT=6020 LEADER_ADDRESS=127.0.0.1:6001 ./subscription-driver &> follower2.txt &disown