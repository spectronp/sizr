#!/bin/sh

# name, isExplicit, size, deps

get_info() {
	pacman_output=$( pacman -Qi "$1" )

	# Name
	echo "$pacman_output" | grep '^Name' | tr -s ' ' | cut -d ' ' -f3

	# Is Explicit
	install_reason=$( echo "$pacman_output" | grep '^Install Reason' | tr -s ' ' | cut -d ' ' -f4-)
	[ "$install_reason" = "Explicitly installed" ] && echo "true" || echo "false"

	# Size
	size_num=$( echo "$pacman_output" | grep '^Installed Size' | tr -s ' ' | cut -d ' ' -f4 )
	size_unit=$( echo "$pacman_output" | grep '^Installed Size' | tr -s ' ' | cut -d ' ' -f5)

	if [ "$size_unit" = "B" ] # TODO -- use elif
	then
		printf "%.0f\n" "$size_num"
	elif [ "$size_unit" = "KiB" ]
	then
		printf "%.0f" "$size_num" | awk '{ print $1 * 1024 }'
	elif [ "$size_unit" = "MiB" ]
	then
		printf "%.0f" "$size_num" | awk '{ print $1 * 1024 ^ 2 }'
	elif [ "$size_unit" = "GiB" ]
	then
		printf "%.0f" "$size_num" | awk '{ print $1 * 1024 ^ 3 }'
	fi

	# Dependencies
	echo "$pacman_output" | grep '^Depends On' | tr -s ' ' | cut -d ' ' -f4-
}

if [ "${0##*/}" = 'info.sh' ]; then
	get_info "$1"
fi
