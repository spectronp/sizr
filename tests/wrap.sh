#!/bin/sh

BASEDIR=$( cd "$(dirname "$0")"; cd ..; pwd ) # TODO -- see if I can do this better
export BASEDIR # NOTE -- do I need this export if i use en to run go test ???
cd "$BASEDIR"

# put pm_mock.sh into a bin dir for each package manager
mkdir bin
cp tests/pm_mock.sh "bin/pacman"
# set PATH
PATH="$BASEDIR/bin:$PATH"
# run go test
env BASEDIR="$BASEDIR" go test "$@"

rm -rf bin
