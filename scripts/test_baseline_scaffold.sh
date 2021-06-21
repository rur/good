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

go run . scaffold baseline/scaffold_test

echo "---- run new server and ping /example endpoint ---"
go run ./baseline/scaffold_test > testing_stdout.log 2> testing_stderr.log &
serverPID=$!
function killserver() {
    echo "Kill test server at PID $serverPID"
    pkill -P $serverPID
}
trap killserver EXIT
sleep 1 # plenty of time to start up
curl http://localhost:8000/example

echo
echo "---- Fetched example page successfully ---"

echo
curl --fail http://localhost:8000/public/test.txt
echo
curl --fail http://localhost:8000/styles/app.css
echo
curl --fail http://localhost:8000/js/app.js

echo
echo "---- Fetched example static files successfully ---"


rm -rf _baseline/scaffold_test
mv baseline/scaffold_test _baseline/
rm -rf baseline

diff=$(git diff _baseline/scaffold_test)

if [[ ! -z $diff ]]; then
    echo "WARNING: Check baseline"
    echo ">>> git diff out >>>"
    printf "$diff"
    echo
    echo ">>> git diff end >>>"
    exit 1
fi

echo "[ok] Scaffold baseline matches!"
