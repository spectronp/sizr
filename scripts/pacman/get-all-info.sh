#!/bin/sh

. ./info.sh

for pack in $( ./get-all.sh ); do
	info=$( get_info "$pack" ) 
	name=$( echo "$info" | sed '1!d' ) # TODO -- is there a better way to do this in POSIX ?
	is_explicit=$( echo "$info" | sed '2!d' )
	size=$( echo "$info" | sed  '3!d' )
	deps=$( echo "$info" | sed '4!d' | sed 's/ /, /g' )

	jq -cnaM \
	--arg name "$name" \
	--arg is_explicit "$is_explicit" \
	--arg size "$size" \
	--arg deps "$deps" \
	'{"name":$name,"isExplicit":$is_explicit,"size":$size,"deps":[$deps]}'
done
