package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type app struct {
	templateRoot string

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	openProject func(string) error
}

func newApp(
	templateRoot string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) *app {
	app := &app{
		templateRoot: templateRoot,
		stdin:        stdin,
		stdout:       stdout,
		stderr:       stderr,
	}

	app.openProject = app.openInCode

	return app
}

func defaultTemplateRoot() (string, error) {
	configRoot, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(
			"get user config directory: %w",
			err,
		)
	}

	return filepath.Join(
		configRoot,
		"mkproj",
		"templates",
	), nil
}

func (a *app) run(args []string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}

	if len(args) == 0 || opts.help {
		showHelp(a.stdout, a.templateRoot)
		return nil
	}

	if opts.list {
		return a.showTemplates()
	}

	if strings.TrimSpace(opts.template) == "" {
		return fmt.Errorf(
			"missing template; use -t <name> or --template=<name>",
		)
	}

	if strings.TrimSpace(opts.target) == "" {
		return fmt.Errorf("missing project path")
	}

	return a.newProject(
		opts.template,
		opts.target,
		opts.open,
	)
}

func (a *app) showTemplates() error {
	entries, err := os.ReadDir(a.templateRoot)

	if os.IsNotExist(err) {
		fmt.Fprintln(a.stdout, "No templates found.")
		fmt.Fprintf(
			a.stdout,
			"Template directory: %s\n",
			a.templateRoot,
		)
		return nil
	}

	if err != nil {
		return err
	}

	var names []string

	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}

	if len(names) == 0 {
		fmt.Fprintln(a.stdout, "No templates found.")
		return nil
	}

	sort.Strings(names)

	for _, name := range names {
		fmt.Fprintln(a.stdout, name)
	}

	return nil
}

func (a *app) templatePath(name string) (string, error) {
	if err := validateTemplateName(name); err != nil {
		return "", err
	}

	path := filepath.Join(a.templateRoot, name)

	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return "", fmt.Errorf(
			"template %q not found",
			name,
		)
	}

	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", fmt.Errorf(
			"template %q is not a directory",
			name,
		)
	}

	return path, nil
}

func validateTemplateName(name string) error {
	if name == "" ||
		name == "." ||
		name == ".." ||
		strings.ContainsAny(name, `/\`) {

		return fmt.Errorf(
			"invalid template name: %q",
			name,
		)
	}

	return nil
}

func resolveTargetPath(target string) (string, error) {
	path, err := filepath.Abs(target)
	if err != nil {
		return "", fmt.Errorf(
			"resolve target path: %w",
			err,
		)
	}

	return filepath.Clean(path), nil
}

func prepareTargetDirectory(path string) error {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0o755)
	}

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf(
			"target exists and is not a directory: %s",
			path,
		)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if len(entries) != 0 {
		return fmt.Errorf(
			"target directory is not empty: %s",
			path,
		)
	}

	return nil
}

func (a *app) newProject(
	template string,
	target string,
	open bool,
) error {
	templatePath, err := a.templatePath(template)
	if err != nil {
		return err
	}

	targetPath, err := resolveTargetPath(target)
	if err != nil {
		return err
	}

	if err := prepareTargetDirectory(targetPath); err != nil {
		return err
	}

	projectName := filepath.Base(targetPath)

	fmt.Fprintf(
		a.stdout,
		"Creating %q from template %q...\n",
		projectName,
		template,
	)

	fmt.Fprintf(a.stdout, "  %s\n", targetPath)

	if err := copyTemplate(
		templatePath,
		targetPath,
	); err != nil {
		return fmt.Errorf(
			"copy template: %w",
			err,
		)
	}

	hookErr := a.invokeTemplateHook(
		targetPath,
		template,
		projectName,
	)

	cleanupErr := os.RemoveAll(
		filepath.Join(
			targetPath,
			".mkproj",
		),
	)

	if hookErr != nil {
		if cleanupErr != nil {
			return fmt.Errorf(
				"%w; cleanup template metadata: %v",
				hookErr,
				cleanupErr,
			)
		}

		return hookErr
	}

	if cleanupErr != nil {
		return fmt.Errorf(
			"cleanup template metadata: %w",
			cleanupErr,
		)
	}

	fmt.Fprintf(
		a.stdout,
		"Created %q.\n",
		projectName,
	)

	if open {
		fmt.Fprintln(a.stdout, "Opening Code...")

		if err := a.openProject(targetPath); err != nil {
			return err
		}
	}

	return nil
}

func copyTemplate(source, target string) error {
	return filepath.WalkDir(
		source,
		func(
			path string,
			entry fs.DirEntry,
			walkErr error,
		) error {
			if walkErr != nil {
				return walkErr
			}

			relative, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}

			if relative == "." {
				return nil
			}

			destination := filepath.Join(
				target,
				relative,
			)

			info, err := entry.Info()
			if err != nil {
				return err
			}

			switch {
			case entry.IsDir():
				return os.MkdirAll(
					destination,
					info.Mode().Perm(),
				)

			case entry.Type()&os.ModeSymlink != 0:
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}

				return os.Symlink(
					link,
					destination,
				)

			case info.Mode().IsRegular():
				return copyFile(
					path,
					destination,
					info.Mode().Perm(),
				)

			default:
				return fmt.Errorf(
					"unsupported file type: %s",
					path,
				)
			}
		},
	)
}

func copyFile(
	source string,
	target string,
	mode fs.FileMode,
) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.OpenFile(
		target,
		os.O_CREATE|
			os.O_TRUNC|
			os.O_WRONLY,
		mode,
	)

	if err != nil {
		return err
	}

	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()

	if copyErr != nil {
		return copyErr
	}

	return closeErr
}

func (a *app) openInCode(path string) error {
	code, err := exec.LookPath("code")
	if err != nil {
		return fmt.Errorf(
			"'code' was not found in PATH",
		)
	}

	cmd := exec.Command(code, path)

	cmd.Stdin = a.stdin
	cmd.Stdout = a.stdout
	cmd.Stderr = a.stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"failed to open Code: %w",
			err,
		)
	}

	return nil
}
