echo "Cleaning up old logs"

rm -rf logs &> /dev/null
mkdir logs

echo "Building the binary"

go build ./

echo "Binary built, executing procs..."

SUB_PORT=6000 LEADER_ADDRESS=0 ./subscription-driver &> ./logs/leader.txt &disown
SUB_PORT=6010 LEADER_ADDRESS=127.0.0.1:6001 ./subscription-driver &> ./logs/follower1.txt &disown
SUB_PORT=6020 LEADER_ADDRESS=127.0.0.1:6001 ./subscription-driver &> ./logs/follower2.txt &disown

echo "Procs are running"