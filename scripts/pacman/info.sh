#!/bin/sh

# name, isExplicit, size, deps

pacman_output=$( pacman -Qi "$1" )

# Name
name=$( echo "$pacman_output" | grep '^Name' | tr -s ' ' | cut -d ' ' -f3 )

# Is Explicit
install_reason=$( echo "$pacman_output" | grep '^Install Reason' | tr -s ' ' | cut -d ' ' -f4-)
[ "$install_reason" = "Explicitly installed" ] && is_explicit="true" || is_explicit="false"

# Version
version=$( echo "$pacman_output" | grep '^Version' | tr -s ' ' | cut -d ' ' -f3 )

# Size
size_num=$( echo "$pacman_output" | grep '^Installed Size' | tr -s ' ' | cut -d ' ' -f4 )
size_unit=$( echo "$pacman_output" | grep '^Installed Size' | tr -s ' ' | cut -d ' ' -f5)

if [ "$size_unit" = "B" ]
then
	size=$( printf "%.0f\n" "$size_num" )	
elif [ "$size_unit" = "KiB" ]
then
	size=$( printf "%.0f" "$size_num" | awk '{ print $1 * 1024 }' )
elif [ "$size_unit" = "MiB" ]
then
	size=$( printf "%.0f" "$size_num" | awk '{ print $1 * 1024 ^ 2 }' )
elif [ "$size_unit" = "GiB" ]
then
	size=$( printf "%.0f" "$size_num" | awk '{ print $1 * 1024 ^ 3 }' )
fi

# Dependencies
deps=$( echo "$pacman_output" | grep '^Depends On' | tr -s ' ' | cut -d ' ' -f4- | sed 's/ /", "/g' )

jq -rcnaM \
	--arg name "$name" \
	--argjson is_explicit "$is_explicit" \
	--arg version "$version" \
	--argjson size "$size" \
	--argjson deps "[\"$deps\"]" \
	'{"name":$name,"isExplicit":$is_explicit,"version": $version,"size":$size,"deps":$deps}'
