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

go run . scaffold baseline/routes_test example trivial

rm baseline/routes_test/page/example/routemap.toml
cp _baseline/testfixtures/routemap.toml baseline/routes_test/page/example/routemap.toml
rm baseline/routes_test/page/trivial/routemap.toml
cp _baseline/testfixtures/routemap_trivial.toml baseline/routes_test/page/trivial/routemap.toml

go run . routes ./baseline/routes_test/page/example
go run . routes ./baseline/routes_test/page/trivial

echo "---- run new server and ping /example endpoint ---"
go run ./baseline/routes_test > testing_stdout.log 2> testing_stderr.log &
serverPID=$!
function killserver() {
    echo "Kill test server at PID $serverPID"
    pkill -P $serverPID
}
trap killserver EXIT
sleep 1 # plenty of time to start up
curl --fail http://localhost:8000/trivial
echo
echo "---"
curl --fail http://localhost:8000/example
echo
echo "---"
curl --fail http://localhost:8000/example/alt
echo
echo "---"
curl --fail http://localhost:8000/example/settings
echo
echo "---"
curl --fail http://localhost:8000/example/advanced-settings
echo
echo "---"
curl --fail -X POST -H "Accept: application/x.treetop-html-template+xml" http://localhost:8000/example/form
echo
echo "---- Feched example page successfully ---"
curl --fail -X POST -H "Accept: application/x.treetop-html-template+xml" http://localhost:8000/example/advanced-settings/submit
echo
echo "---- Feched example page successfully ---"
curl -X POST -H "Accept: application/x.treetop-html-template+xml" http://localhost:8000/example/advanced-settings/submit
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
    printf "$diff"
    echo
    echo ">>> git diff end >>>"
    exit 1
fi

echo "[ok] Scaffold baseline matches!"
