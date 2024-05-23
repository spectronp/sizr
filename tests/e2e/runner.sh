#!/usr/bin/env bash

package_managers=('apt' 'dnf' 'pacman') # TODO -- get these from an env var ?

apk update &&
    apk upgrade &&
    apk install bash

# NOTE -- is this needed ?
# hide Alpine package manager
mv /usr/bin/apk /usr/bin/.apk

# TODO -- just change the PATH
for package_manager in "${package_managers[@]}"; do
    mv pm_mock.sh /usr/bin/${package_manager}

    venom run
done
