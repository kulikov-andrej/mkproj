# mkproj

`mkproj` is a small cross-platform project generator built around reusable directory templates, optional Starlark setup scripts, and per-project workflows.

## Installation
Prebuilt releases currently target Windows and Linux on amd64, and macOS on amd64 and arm64.

### Windows

Download and run the installer from the latest release:

```powershell
$uri = "https://github.com/kulikov-andrej/mkproj/releases/latest/download/mkproj-install-windows-amd64.exe"
$installer = Join-Path $env:TEMP "mkproj-install.exe"
Invoke-WebRequest $uri -OutFile $installer
& $installer
```

The installer places `mkproj.exe` under `%LOCALAPPDATA%\Programs\mkproj` and adds that directory to the user `PATH` if necessary. Open a new terminal after installation when prompted.

### Linux

```sh
curl -fL https://github.com/kulikov-andrej/mkproj/releases/latest/download/mkproj-install-linux-amd64 -o /tmp/mkproj-install
chmod +x /tmp/mkproj-install
/tmp/mkproj-install
```

The installer places `mkproj` in `~/.local/bin`. If that directory is not already in `PATH`, the installer prints a reminder.

### macOS

The installer is available for both Apple Silicon and Intel Macs:

```sh
case "$(uname -m)" in
  arm64) arch=arm64 ;;
  x86_64) arch=amd64 ;;
  *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

curl -fL "https://github.com/kulikov-andrej/mkproj/releases/latest/download/mkproj-install-darwin-$arch" -o /tmp/mkproj-install
chmod +x /tmp/mkproj-install
/tmp/mkproj-install
```

The installer places mkproj in ~/.local/bin. If that directory is not already in PATH, the installer prints a reminder.

## Quick start

List the templates available to the current user:

```console
mkproj template list
```

Create a project from the `cpp` template:

```console
mkproj hello-world -t cpp
```

Create a project in the current directory:

```console
mkproj -t cpp
```

Create a project and open it in Visual Studio Code:

```console
mkproj hello-world -t cpp --open
```

If the generated project defines a workflow, list its commands with:

```console
mkproj run
```

and run one with:

```console
mkproj run build
```

## Commands

```text
mkproj [<path>] -t <template> [--open]
mkproj template list
mkproj run [<command>]
mkproj version
mkproj help [<topic>]
```

`<path>` defaults to the current directory. The target directory may be absent or already exist and be empty; `mkproj` does not overwrite a non-empty directory.

`--open` (or `-o`) opens the generated project with the `code` executable after creation.

## Templates

A template is an ordinary directory stored in the user's configuration directory. It may also contain `.mkproj/setup.star`, which runs after the template is copied and can customize the generated project.

See [Templates](docs/templates.md) for the directory layout, setup lifecycle, and Starlark API.

## Project workflows

A generated project may keep `.mkproj/workflow.star` and expose reusable project commands through `mkproj run`.

See [Workflows](docs/workflows.md) for the workflow format and examples.

## License

See [LICENSE](LICENSE).
