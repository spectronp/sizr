#!/usr/bin/env bash

# name type size deps...

# TODO -- deal with circular dependencies
# TODO -- generate variables closer to json output

while [[ "$#" -gt 0 ]]
do
    case "$1" in
        '-n')
            exp_packages_num="$2"
            shift 2
            ;;
        '-d')
            fixed_deps_num="$2"
            shift 2
            ;;
        '--odd')
            new_dep_odd="$2"
            shift 2
            ;;
    esac
done

names=""
gen_name() {
    until [[ ! " ${names[*]} " =~ " $name " ]]
    do
        name=$(tr -dc 'a-z' </dev/urandom | head -c5)
    done

    names+=" $name"

    echo "$name"
}

gen_size() {
    echo $(($RANDOM * 100))
}

gen_deps_num() {
    [[ -z "$fixed_deps_num" ]] && echo $(( $RANDOM % 6 )) || echo "$fixed_deps_num"
}

odd() {
    num_pick=$(( $RANDOM % 100 + 1 ))

    [[ "$num_pick" -le $1 ]] && return 0
}

output() { # TODO -- fix: the , after the last package
    echo "$*" >> "$output_file"

    [ "$2" == "exp" ] && is_exp=true || is_exp=false
    
    deps=$( echo "$4" | sed -e 's/,/", "/g' )

    cat <<EOF >> packages.json
            "$1": {
               "name": "$1",
               "isExplicit": $is_exp,
               "size": $3,
               "deps": [ "$deps" ] 
            },
EOF
}

cat <<EOF > packages.json
{
    "manager": "pm",
    {
EOF

output_file="/dev/stdout"
# output_file="packages-gen.txt"

echo '' > "$output_file"
deps=()

# create primary packages

[[ -z "$exp_packages_num" ]] && exp_packages_num=$(( 15 + $RANDOM % 6 ))

for ((i=1; i <= $exp_packages_num; i++)); do
    package="`gen_name` exp `gen_size`"
    deps_num=`gen_deps_num`

    [ "${#deps_num[@]}" -eq 0 ] && continue
    package+=' '

    for ((i2=1; i2 <= $deps_num; i2++)); do
        dep_name=`gen_name`
        deps+=("$dep_name")
        package+="$dep_name,"
    done
    output ${package::-1}
done

# create dependencies packages

deps_not_created=("${deps[@]}")
[[ -z "$new_dep_odd" ]] && new_dep_odd=30
new_dep_odd_decrease=3
while [[ "${#deps_not_created[@]}" -gt 0 ]]; do
    dep="${deps_not_created[0]}"
    deps_not_created=("${deps_not_created[@]:1}")

    line="$dep dep $(gen_size)"
    deps_num=`gen_deps_num`

    [[ "$deps_num" -eq 0 ]] && output $line && continue
    line+=' '

    for ((i=1; i<=$deps_num; i++)); do
        if odd $new_dep_odd; then
            let "new_dep_odd-=$new_dep_odd_decrease"
            new_dep=$(gen_name)
            deps_not_created+=("$new_dep")
            line+="$new_dep,"
        else
            while true
            do
                dep_index=$(( $RANDOM % ${#deps[@]} ))
                dep_name="${deps[$dep_index]}"
                [[ "$line" =~ "$dep_name" ]] || line+="$dep_name," && break
            done
        fi
    done
    output ${line::-1}
done

echo "}" >> packages.json
