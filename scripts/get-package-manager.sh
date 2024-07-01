#!/bin/sh

# try apt
command -v apt &> /dev/null && echo -n "apt" && exit 0

# try pacman
command -v pacman &> /dev/null && echo -n "pacman" && exit 0

exit 1
