#!/bin/sh

# try apt
apt -v >/dev/null 2>&1 && echo -n "apt" && exit 0

# try pacman
pacman --version >/dev/null 2>&1 && echo -n "pacman" && exit 0

exit 1
