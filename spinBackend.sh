# Clean up old logs
rm -rf ./logs &>/dev/null
mkdir logs
touch logs/blurb-log.txt
touch logs/sub-log.txt
touch logs/reg-log.txt

# Spawn backend services, piped to log files
cd blurb-driver && ./run.sh &>../logs/blurb-log.txt &disown
cd subscription-driver && ./run.sh &>../logs/sub-log.txt &disown
cd registration-driver && ./run.sh &>../logs/reg-log.txt &disown

echo "Backend services started..."
