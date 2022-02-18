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

go run . scaffold baseline/routes_test

go run . page ./baseline/routes_test example --starter :minimum -y
rm baseline/routes_test/page/example/routemap.toml
cp _baseline/testfixtures/routemap.toml baseline/routes_test/page/example/routemap.toml

go run . page ./baseline/routes_test trivial --starter :minimum -y --no-resources
rm baseline/routes_test/page/trivial/routemap.toml
cp _baseline/testfixtures/routemap_trivial.toml baseline/routes_test/page/trivial/routemap.toml

go run . routes gen ./baseline/routes_test/page/example
go run . routes gen ./baseline/routes_test/page/trivial

if [[ ! -z $(bash ./scripts/usedports.sh | grep 8001) ]]; then
  echo >&2 "Port 8001 appears to be in use, cannot run test"
  exit 1
fi

echo "---- run new server and ping /example endpoint ---"
go run ./baseline/routes_test --port 8001 > testing_stdout.log 2> testing_stderr.log &
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
curl --fail http://localhost:8001/trivial > _test_output/1.html
echo
echo "---"
curl --fail http://localhost:8001/example > _test_output/2.html
echo
echo "---"
curl --fail http://localhost:8001/example/alt > _test_output/3.html
echo
echo "---"
curl --fail -H "Accept: application/x.treetop-html-template+xml" http://localhost:8001/example/settings > _test_output/4.html
echo
echo "---"
curl --fail http://localhost:8001/example/advanced-settings > _test_output/5.html
echo
echo "---"
curl --fail -X POST -H "Accept: application/x.treetop-html-template+xml" http://localhost:8001/example/form > _test_output/6.html
echo
curl --fail -X POST -H "Accept: application/x.treetop-html-template+xml" http://localhost:8001/example/advanced-settings/submit > _test_output/7.html
echo
echo "---- Feched example page successfully ---"


rm -rf _baseline/routes_test
mv baseline/routes_test _baseline/
rm -r baseline

# normalize name of generated handlers file for comparison against baseline
mv _baseline/routes_test/page/example/handlers_* _baseline/routes_test/page/example/handlers_gen.go
mv _baseline/routes_test/page/trivial/handlers_* _baseline/routes_test/page/trivial/handlers_gen.go

diff=$(git diff _baseline/routes_test)

if [[ ! -z $diff ]]; then
    echo "WARNING: Check baseline"
    echo ">>> git diff out >>>"
    echo "$diff"
    echo
    echo ">>> git diff end >>>"
    exit 1
fi

echo "[ok] Scaffold baseline matches!"
