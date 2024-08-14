#!/usr/bin/env bash

packages_file="$BASEDIR/tests/packages.json"

cmd_types=('list' 'info')

# list -> list explictly installed packages
# info -> display info of 1 specified package

apt_list='apt-mark showmanual'

dpkg_info="dpkg-query -f \${Package} \${Installed-Size} \${Pre-Depends},\${Depends} -W"

pacman_list='pacman -Q'
pacman_info='pacman -Qi'

package_manager=$( echo "${0##*/}" | grep -Eo '^[a-z]+' )

to_human_readable() {
    raw_size=$1
    shift
    shift

    prefixes=('B' 'K' 'M' 'G')
    base=1024

    declare prefix_num
    for i in {3..1}
    do
        (( $raw_size / $base ** $i )) && prefix_num="$i" && break # TODO: check the error here
    done

    prefix_num="${prefix_num:-0}"
    size=$(printf %.2f "$(( $raw_size / $base ** $prefix_num ))") # TODO: numbers are being rounded ( e.g 5 / 2 = 2.00 )
    prefix=${prefixes[$prefix_num]}
    unit="${prefix}iB"

    echo "$size $unit"
}

apt_output() {
    cmd_type="$1"
    case "$cmd_type" in
        'list')
            jq -r '.[] | if .isExplicit then . else empty end' "$packages_file"
            ;;
    esac
}

dpkg_output() {
    cmd_type="$1"
    case "$cmd_type" in
        'info')
            package_name="$5" # TODO: automate it
            cat "$packages_file" | awk "\$1 == \"$package_name\" { print \$1,\$3,\$4 }" # TODO: use jq here
            ;;
    esac
}

pacman_output() {
    cmd_type="$1"
    case "$cmd_type" in
        'list')
            jq -r '.[] | .name as $name | .version as $version | "\($name) \($version)"' "$packages_file"
            ;;
        'info')
            package_name="$3"
            package_size="10 KiB"
            # package_size=$(jq -r ".$package_name.size" "$packages_file")
            # package_size=$(to_human_readable $package_size)

            # TODO: make this more readable
            deps=$(jq -r "reduce .$package_name.deps[] as \$dep (\"\"; . + \" \" + \$dep )" "$packages_file")
            package_type=$(jq -r ".$package_name.isExplicit" "$packages_file")
            [[ "$package_type" == "true" ]] && package_type="Explicitly installed" || package_type="Installed as a dependency for another package"
            cat "$BASEDIR/tests/outputs/pacman_info.txt" | sed "s/{name}/$package_name/;s/{deps}/$deps/;s/{size}/$package_size/;s/{package_type}/$package_type/"
            ;;
    esac
}

declare cmd
for cmd_type in "${cmd_types[@]}"
do
    cmd_ref="${package_manager}_"
    cmd_ref+="$cmd_type"

    if [[ "$cmd_type" == 'info' ]]
    then
        args=$( echo "$@" | sed -E 's/ [[:alnum:]]+//' )
    else
        args="$@"
    fi

    if [[ $( echo "${0##*/}" "$args" ) == "${!cmd_ref}" ]]; then
        cmd="$cmd_type"
        break
    fi
done

# TODO: error if $cmd is empty

${package_manager}_output $cmd "$@"
