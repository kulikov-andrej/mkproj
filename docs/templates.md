# Templates

A template is an ordinary directory that `mkproj` copies into a new project.

## Location

Templates are stored under the current user's configuration directory:

```text
<user-config>/mkproj/templates/
```

Typical locations are:

```text
Windows: %AppData%\mkproj\templates\
Linux:   $XDG_CONFIG_HOME/mkproj/templates/
         or ~/.config/mkproj/templates/
macOS:   ~/Library/Application Support/mkproj/templates/
```

Each direct child directory is one template:

```text
templates/
├── cpp/
├── go/
└── python/
```

The directory name is the template name:

```console
mkproj hello-world -t cpp
```

Use `mkproj template list` to list available templates.

## Starter template

For a quick first run, `mkproj` can install a small example template:

```console
mkproj template init
```

It creates only the `starter` template and leaves any existing templates untouched:

```text
templates/
└── starter/
    ├── .mkproj/
    │   └── setup.star
    ├── README.md
    └── src/
        └── main.txt
```

The setup script replaces `{{PROJECT_NAME}}` in the example files. Because `setup.star` is removed after setup, a generated project has a simpler layout:

```text
first-project/
├── README.md
└── src/
    └── main.txt
```

Try the complete flow with:

```console
mkproj template init
mkproj template list
mkproj first-project -t starter -o
mkproj template deinit
```

`mkproj template init` creates `starter`; if it is already installed, the command reports the existing path and exits successfully without overwriting it. `mkproj template deinit` removes only `starter`; if it is already absent, the command reports that nothing is installed and exits successfully.

## Project creation

A template may contain any files and directories:

```text
cpp/
├── CMakeLists.txt
├── include/
└── src/
    └── main.cpp
```

The complete template tree is copied into the target directory. The target may not exist, or it may already exist and be empty. A file path or a non-empty directory is rejected.

If no target path is given, `mkproj` uses the current directory:

```console
mkproj -t cpp
```

## Setup script

A template may optionally contain:

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

The setup script runs after the template has been copied. It can inspect project and template information and modify files inside the generated project.

```python
print("Configuring " + project.name)

replace(
    "README.md",
    "{{PROJECT_NAME}}",
    project.name,
)

run("git", "init")
```

The following values are available:

```text
project.name    Generated project directory name
project.path    Absolute path to the generated project
template.name   Template name passed to mkproj
```

### Built-ins

`run(program, *args)` executes a program directly, without a shell, with the generated project as its working directory. Arguments must be strings. Standard input, output, and error are connected to `mkproj`.

```python
run("git", "init")
run("go", "mod", "tidy")
```

`replace(path, old, new)` replaces every occurrence of `old` in a file and returns the number of replacements. Existing file permissions are preserved.

```python
count = replace("README.md", "{{NAME}}", project.name)
```

`write(path, content)` writes a file, creating parent directories when necessary.

```python
write("config/generated.txt", project.name + "\n")
```

`mkdir(path)` creates a directory and any missing parents.

```python
mkdir("src/generated")
```

`remove(path)` removes a file or directory recursively. Removing the project root itself is refused.

```python
remove("example")
```

The path-based built-ins use paths relative to the generated project. Empty paths, absolute paths, and paths that escape the project with `..` are rejected.

Starlark `load()` is not supported. Setup scripts are trusted code: `run()` can execute programs available to the current user, so templates should be treated like any other executable project bootstrap code.

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
remove .mkproj/setup.star
      ↓
remove .mkproj if it is empty
      ↓
finish
```

Cleanup is attempted even when setup fails. A failed setup does not roll back files already copied or changes already made by the script.

If the template contains other `.mkproj` metadata, such as `workflow.star`, that metadata remains in the generated project.
