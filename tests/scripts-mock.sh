#!/usr/bin/env bash

BASE_PATH=$( dirname "$0" )

script="$1"
shift

case "$script" in
	"get-package-manager")
		echo -n "pm"
		;;
	"pm/get-all")
		jq -r '.[] | .name' "$BASE_PATH/packages.json"
		;;
	"pm/info")
		jq -caM ".$1" "$BASE_PATH/packages.json"

esac
