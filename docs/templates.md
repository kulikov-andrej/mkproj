# Templates

A template is an ordinary directory that `mkproj` copies into a new project.

## Location

Templates are stored in the current user's configuration directory:

```text
<user-config>/mkproj/templates/
```

Typical locations:

```text
Windows: %AppData%\mkproj\templates\
Linux:  $XDG_CONFIG_HOME/mkproj/templates/
        or ~/.config/mkproj/templates/
```

Each direct child directory is a template:

```text
templates/
├── cpp/
├── go/
└── python/
```

The directory name is used as the template name:

```console
mkproj hello_world -t cpp
```

## Structure

A template may contain any files and directories:

```text
cpp/
├── CMakeLists.txt
├── include/
└── src/
    └── main.cpp
```

When a project is created, the template is copied into the target directory.

The target may not exist or may already exist and be empty. `mkproj` does not overwrite a non-empty directory.

## Setup

Templates may optionally contain a setup script:

```text
.mkproj/setup.star
```

For example:

```text
python/
├── .mkproj/
│   └── setup.star
├── README.md
├── main.py
└── pyproject.toml
```

The setup script runs after the template has been copied.

Example:

```starlark
print("Configuring " + project.name)

replace(
    "README.md",
    "{{PROJECT_NAME}}",
    project.name,
)

run("git", "init")
```

Setup scripts have access to project and template information:

```text
project.name
project.path
template.name
```

The following functions are available:

```text
run(...)
replace(...)
write(...)
mkdir(...)
remove(...)
```

The path-based functions `replace()`, `write()`, `mkdir()`, and `remove()` use paths relative to the generated project. Absolute paths and paths that escape the project directory with `..` are rejected.

`run()` executes programs directly, without a shell, with the generated project as the working directory.

Setup scripts are trusted code, not a security sandbox.

## Lifecycle

Project creation follows this order:

```text
resolve template
      ↓
validate target
      ↓
copy template
      ↓
run .mkproj/setup.star if present
      ↓
remove .mkproj
      ↓
finish
```

The `.mkproj` directory is removed before project creation finishes. Cleanup is attempted even if setup fails.

A failed setup does not roll back changes already made to the project.
