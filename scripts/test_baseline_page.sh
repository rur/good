#!/usr/bin/env bash

set -e

# sanity check
if [[ ! $(head -n 1 ./go.mod) == "module github.com/rur/good" ]]; then
    echo "This script must be run from github.com/rur/good project directory"
    exit 1
fi

cat <<TESTINFO
--- testing good page command ---
TESTINFO

echo "clearing any previously failed run data"
rm -rf baseline

go run . scaffold baseline/page_test
go run . page baseline/page_test newpage

echo "---- run new server and ping endpoints ---"
go run ./baseline/page_test > testing_stdout.log 2> testing_stderr.log &
serverPID=$!
function killserver() {
    echo "Kill test server at PID $serverPID"
    pkill -P $serverPID
}
trap killserver EXIT
sleep 1 # plenty of time to start up

echo
curl --fail http://localhost:8000/intro
echo
curl --fail http://localhost:8000/newpage

echo "---- Feched intro and newpage page successfully ---"

rm -rf _baseline/page_test
mv baseline/page_test _baseline/
rm -r baseline

diff=$(git diff _baseline/page_test)

if [[ ! -z $diff ]]; then
    echo "WARNING: Check baseline"
    echo ">>> git diff out >>>"
    printf "$diff"
    echo
    echo ">>> git diff end >>>"
    exit 1
fi

echo "[ok] Page scaffold baseline matches!"