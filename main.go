package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tj/docopt"
	"github.com/ttacon/chalk"
)

const version = "0.3.0"

const usage = `
Usage:
  master [options] <file-or-directory>
  master -h | --help
  master --version

Options:
  -d, --output-directory string  Specify the output directory (default: <file-or-directory>).
  -s, --schema-directory string  Specify the JSON Schema directory (default: <file-or-directory>).
  -e, --encoding string          CSV file encoding [default: auto]. Supported encodings are https://goo.gl/T3zICN
  -E, --fix-encoding             Fix the CSV file encoding if it is different from --encoding.
  -n, --no-output-file           No file output. If file is given, print JSON string to stdout.
  -S, --output-schema            Output JSON Schema from CSV files.
  -V, --skip-validation          Skip validation by JSON Schema.
  -j, --no-schema-suffix         Disable to use *.schema.json suffix pattern.
  -h, --help                     Output help information.
  -v, --version                  Output version.
`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, true)
	if err != nil {
		fatalf("Failed to parse arguments: %v\n%v", args, err)
	}

	var file, dir string
	fileOrDir := resolvePath(args["<file-or-directory>"].(string))
	stat, err := os.Stat(fileOrDir)
	if os.IsNotExist(err) {
		fatalf("No such file or directory: %v\n%v", fileOrDir, err)
	}
	if stat.IsDir() {
		dir = fileOrDir
	} else {
		file = fileOrDir
		dir = filepath.Dir(file)
	}

	if args["--output-directory"] == nil {
		args["--output-directory"] = dir
	}
	if args["--schema-directory"] == nil {
		args["--schema-directory"] = dir
	}

	cli := &Cli{
		dir:            dir,
		file:           file,
		outputDir:      resolvePath(args["--output-directory"].(string)),
		schemaDir:      resolvePath(args["--schema-directory"].(string)),
		encoding:       args["--encoding"].(string),
		fixEncoding:    args["--fix-encoding"].(bool),
		noOutputFile:   args["--no-output-file"].(bool),
		outputSchema:   args["--output-schema"].(bool),
		skipValidation: args["--skip-validation"].(bool),
		noSchemaSuffix: args["--no-schema-suffix"].(bool),
	}
	cli.run()
}

func fatalf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "\n"+chalk.Red.Color("[Error]")+" %s\n\n", fmt.Sprintf(msg, args...))
	os.Exit(1)
}

func resolvePath(path string) string {
	resolved, err := filepath.Abs(path)
	if err != nil {
		fatalf("Failed to resolve path: %v\n%v", path, err)
	}
	return resolved
}
