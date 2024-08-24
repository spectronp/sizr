# sizr

> A tool that recursively calculates the removable size of packages

It will sum the size of every dependency that will no longer be required if the target package is deleted.

By default, **sizr** will display the packages of higher sizes in decrescent order.

## How to use

The first time you run **sizr**, it will load all the packages installed on your system ( This might take a minute, but only on the first time ).
After, it will display a list ( 30, by default ) of the packages with higher calculated sizes. The number of packages can be changed with `-n` or `--limit` flags.

### Usage

```
sizr [-n | --limit] ( Default: 30 )
sizr -v | --version
sizr -h | --help

Options:
  -n --limit	Set the limit of packages to show (Default: 30)
  -h --help	Show this help message
  -v --version	Show sizr version
```

## How to install

```
git clone https://github.com/spectronp/sizr
cd sizr
make
make install
```

**IMPORTANT!** The user needs to be in the `wheel` group to use the database file

Add current user to `wheel`:

```
sudo usermod -a -G wheel $USER
```

### Dependencies

The only dependency **sizr** has is [jq](https://jqlang.github.io/jq/). You can easily install it with:

```
sudo pacman -S jq

# or using an AUR helper ( I think you know how to do it 🙂 )
```

### How to uninstall

Run inside the repository:

```
make uninstall
```

## How to run tests

Run the default test command:

```
make test
```

You can pass arguments and flags to `go test` command running the `wrap.sh` script:

```
./test/wrap.sh ARGS_FLAGS_HERE

# Example
./test/wrap.sh -v -failfast
```

## Package Managers Supported

- [x] Pacman
- [ ] Apt ( Planned in the near future )

## Roadmap

- [ ] Support more Package Managers
- [ ] Config
- [ ] Overall Optimization
- [ ] More info about packages
- [x] CLI
  - [x] List report
  - [ ] Single package
  - [ ] Merge mode ( calculate more than 1 package )
  - [ ] Utility flags
  - [ ] Log
  - [ ] Benchmark
- [ ] TUI
