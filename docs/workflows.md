# Workflows

A project can define reusable commands in:

```text
.mkproj/workflow.star
```

`mkproj` treats top-level callable values in that file as workflow commands.

## Example

```python
def build():
    run("go", "build", "./...")


def test():
    run("go", "test", "./...")


def clean():
    remove("build")
```

From anywhere inside that project, list the available commands:

```console
mkproj run
```

Example output:

```text
Usage:
  mkproj run <command>

Commands:
  build
  clean
  test
```

Run a command with:

```console
mkproj run test
```

Workflow commands are called without arguments. Commands are listed alphabetically.

## Project discovery

`mkproj run` starts from the current working directory and walks upward until it finds a `.mkproj` directory. This makes workflows available from nested directories inside the project.

If no project metadata is found, running a workflow command fails with `not inside an mkproj project`.

## Starlark API

Workflow scripts receive project information:

```text
project.name    Project directory name
project.path    Absolute path to the project root
```

They use the same built-ins as template setup scripts:

```text
run(...)
replace(...)
write(...)
mkdir(...)
remove(...)
```

`run()` executes a program directly with the project root as its working directory and forwards stdin, stdout, and stderr.

The path-based built-ins operate relative to the project root and reject paths that escape it. `remove()` also refuses to remove the project root itself.

See [Templates](templates.md#built-ins) for the built-in reference.

Starlark `load()` is not supported. Workflow files are trusted project code and may execute external programs through `run()`.
