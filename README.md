# mkproj

`mkproj` is a tiny cross-platform utility for creating projects from reusable templates.

## Installation

### Windows

```powershell
irm https://raw.githubusercontent.com/kulikov-andrej/mkproj/master/install/windows.ps1 | iex
```

### Linux

```sh
curl -fsSL https://raw.githubusercontent.com/kulikov-andrej/mkproj/master/install/linux.sh | bash
```

## Quick start

List templates:

```console
mkproj -l
```

Create a `hello_world` project from the `python` template:

```console
mkproj hello_world -t python
```

Create and open a `cpp_tutorial` project from the `cpp` template:

```console
mkproj cpp_tutorial -t cpp -o
```

Create a project in the current directory:

```console
mkproj -t cpp
```

## Documentation

- [Templates](docs/templates.md)
