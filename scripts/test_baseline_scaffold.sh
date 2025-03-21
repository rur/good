#!/usr/bin/env bash

set -e

# sanity check
if [[ ! $(head -n 1 ./go.mod) == "module github.com/rur/good" ]]; then
    echo "This script must be run from github.com/rur/good project directory"
    exit 1
fi

cat <<TESTINFO
--- testing good scaffold command ---
TESTINFO

echo "clearly any previously failed run data"
rm -rf baseline

go run . scaffold ./baseline/scaffold_test

if [[ ! -z $(bash ./scripts/usedports.sh | grep 8001) ]]; then
  echo >&2 "Port 8001 appears to be in use, cannot run test"
  exit 1
fi

echo "---- run new server and ping /example endpoint ---"
go run ./baseline/scaffold_test --port 8001 > testing_stdout.log 2> testing_stderr.log &
serverPID=$!
function killserver() {
    echo "Kill test server at PID $serverPID"
    pkill -P $serverPID
}
trap killserver EXIT
sleep 1 # plenty of time to start up

rm -rf _test_output
mkdir _test_output
curl http://localhost:8001/example > _test_output/1.html

echo
echo "---- Fetched example page successfully ---"

echo
curl --fail http://localhost:8001/public/test.txt > _test_output/2.html
echo
curl --fail http://localhost:8001/styles/app.css > _test_output/3.html
echo
curl --fail http://localhost:8001/js/app.js > _test_output/4.html

echo
echo "---- Fetched example static files successfully ---"


rm -rf _baseline/scaffold_test
mv baseline/scaffold_test _baseline/
rm -rf baseline

diff=$(git diff _baseline/scaffold_test)

if [[ ! -z $diff ]]; then
    echo "WARNING: Check baseline"
    echo ">>> git diff out >>>"
    echo "$diff"
    echo
    echo ">>> git diff end >>>"
    exit 1
fi

echo "[ok] Scaffold baseline matches!"
