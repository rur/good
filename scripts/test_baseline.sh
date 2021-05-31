#!/usr/bin/env bash

set -e

# sanity check
if [[ ! $(head -n 1 ./go.mod) == "module github.com/rur/good" ]]; then
    echo "This script must be run from github.com/rur/good project directory"
    exit 1
fi

rm -rf _baseline/site/*
go run . scaffold _baseline/site

diff=$(git diff _baseline/site)

if [[ ! -z $diff ]]; then
    echo "WARNING: Check baseline"
    echo ">>> git diff out >>>"
    printf "$diff"
    echo
    echo ">>> git diff end >>>"
    exit 1
fi

echo "Scaffold baseline matches!"

function restoreBackup() {
    exitcode=$?
    if [[ -d _baseline/site_page_bk ]]; then
        # restore site from backup during page test
        rm -r _baseline/page
        mv _baseline/site _baseline/page
        mv _baseline/site_page_bk _baseline/site
    fi
    exit $exitcode
}

trap restoreBackup EXIT

cp -r _baseline/site _baseline/site_page_bk/

# run page generator on the _baseline/site scaffold
go run . page _baseline/site settings

# prepare _baseline/pages working dir for comparison and restore site
rm -rf _baseline/page
mv _baseline/site _baseline/page
mv _baseline/site_page_bk _baseline/site

diff=$(git diff _baseline/page)

if [[ ! -z $diff ]]; then
    echo "WARNING: Check baseline"
    echo ">>> git diff out >>>"
    printf "$diff"
    echo
    echo ">>> git diff end >>>"
    exit 1
fi

echo "Page scaffold baseline matches!"