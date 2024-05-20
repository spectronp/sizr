#!/usr/bin/env bash

script="$1"
shift

case "$script" in
	"get-package-manager")
		echo "pm"
		;;
	"pm/get-all")
		awk '{ print $1 }' ./e2e/packages-gen.txt
		;;
	"pm/info")
		grep "^$1 " ./e2e/packages-gen.txt | sed -e 's/ /\n/g' -e 's/,/ /g' -e 's/exp/true/' -e 's/dep/false/'
esac
