#!/bin/sh

BASEDIR=$( cd "$(dirname "$0")"/..; pwd )
cd "$BASEDIR"

# put pm_mock.sh into a bin dir for each package manager
mkdir bin
cp tests/pm_mock.sh "bin/pacman"
# set PATH
PATH="$BASEDIR/bin:$PATH"
# run go test
env BASEDIR="$BASEDIR" SIZR_ENV=testing go test ./... "$@"

rm -rf bin
