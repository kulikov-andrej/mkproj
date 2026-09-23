# mkproj

`mkproj` is a small cross-platform project generator built around reusable directory templates, optional Starlark setup scripts, and per-project workflows.

## Installation

Prebuilt releases currently target Windows and Linux on amd64, and macOS on amd64 and arm64.

Download and run the installer from the latest release.

## Quick start

Install the example template, create a project, and clean the example up when you are done:

```console
mkproj template init
mkproj template list
mkproj first-project -t starter -o
mkproj template deinit
```

The generated project starts with a small example structure:

```text
first-project/
├── README.md
└── src/
    └── main.txt
```

`template deinit` removes only the bundled `starter` template. To create a project in the current directory instead, omit the path (the directory must be empty):

```console
mkproj -t starter
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
mkproj template [<subcommand>]
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
