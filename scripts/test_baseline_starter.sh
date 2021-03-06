#!/usr/bin/env bash

set -e

# sanity check
if [[ ! $(head -n 1 ./go.mod) == "module github.com/rur/good" ]]; then
    echo "This script must be run from github.com/rur/good project directory"
    exit 1
fi

cat <<TESTINFO
--- testing good starter command ---
TESTINFO

echo "clearing any previously failed run data"
rm -rf baseline

set -x
go run . scaffold baseline/starter_test
go run . starter baseline/starter_test/starter
go run . page ./baseline/starter_test newpage --starter ./baseline/starter_test/starter -y
set +x

if [[ ! -z $(bash ./scripts/usedports.sh | grep 8001) ]]; then
  echo >&2 "Port 8001 appears to be in use, cannot run test"
  exit 1
fi

echo "---- run new server and ping endpoints ---"
go run ./baseline/starter_test --port 8001 > testing_stdout.log 2> testing_stderr.log &
serverPID=$!
function killserver() {
    echo "Kill test server at PID $serverPID"
    pkill -P $serverPID
}
trap killserver EXIT
sleep 1 # plenty of time to start up

rm -rf _test_output
mkdir _test_output

echo
curl --fail http://localhost:8001/example > _test_output/1.html
echo
curl --fail http://localhost:8001/newpage > _test_output/2.html

echo "---- Feched example and newpage page successfully ---"

rm -rf _baseline/starter_test
mv baseline/starter_test _baseline/
rm -r baseline

diff=$(git diff _baseline/starter_test)

if [[ ! -z $diff ]]; then
    echo "WARNING: Check baseline"
    echo ">>> git diff out >>>"
    echo "$diff"
    echo
    echo ">>> git diff end >>>"
    exit 1
fi

echo "[ok] Page scaffold baseline matches!"
