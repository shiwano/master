# master [![Build Status](https://secure.travis-ci.org/shiwano/master.png?branch=master)](http://travis-ci.org/shiwano/master)

Converts CSV to structured JSON with JSON Schema validation.

Example:

|id|author.id|author.name|author.is_human|members.0.id|members.0.name|members.1.id|members.1.name|comments.0|comments.1|
|---|---|---|---|---|---|---|---|---|---|
|1|100|Alice|TRUE|200|White rabbit|300|Cheshire Cat|Hi|Hello|

master will convert this CSV to like below:

```json
[
  {
    "id": 1,
    "author": { "id": 100, "name": "Alice", "is_human": true },
    "members": [
      { "id": 200, "name": "White Rabbit" },
      { "id": 300, "name": "Cheshire Cat" }
    ],
    "comments": [ "Hi", "Hello" ]
  }
]
```

## Installation

Via binary [releases](https://github.com/shiwano/master/releases).

Via `go-get`:

```bash
$ go get -u github.com/shiwano/master
```

Via [Homebrew](http://brew.sh/):

```bash
brew tap shiwano/formulas
brew install master
```

## Usage

```
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
  -S, --output-schema            Output JSON schema from CSV files.
  -V, --skip-validation          Skip validation by JSON Schema.
  -j, --no-schema-suffix         Disable to use *.schema.json suffix pattern.
  -h, --help                     Output help information.
  -v, --version                  Output version.
```

## Nested Object and Array

master uses dot(.) as a separator to clarify nested object and array.
You can use csv column name patterns like below.

Nested object:

|user.id|user.name|
|---|---|
|100|Alice|

```json
[
  { "user": { "id": 100, "name": "Alice" } }
]
```

Array:

|items.0|items.1|items.2|
|---|---|---|
|1|2|3|

```json
[
  { "items": [ 1, 2, 3 ] }
]
```

Mix:

|users.0.id|users.0.name|users.1.id|users.1.name|
|---|---|---|---|
|100|Alice|200|White Rabbit|

```json
[
  {
    "users": [
      { "id": 100, "name": "Alice" },
      { "id": 200, "name": "White Rabbit" }
    ]
  }
]
```

## Validation

master supports JSON Schema validation. For example,
if `foo.csv` was given as argument, master finds `foo.schema.json` from
schema directory, and will use it for validation.

The `--output-schema` option lets you get easily JSON Schema from CSV.

```bash
$ master --output-schema masterdata.csv
```

## Encoding

master uses [chardet](https://github.com/saintfish/chardet) libraly to detect
charset of CSV files. It's based on the algorithm and data in
[ICU](http://icu-project.org/)'s implementation.

If the `--fix-encoding` option was given, master fixes the CSV file encoding
to the `--encoding` option value (`auto` is same as `UTF-8`).
Note that master always adds BOM to UTF-8 CSV files.

## Boolean

master parses CSV's `TRUE` and `FALSE` strings to JSON's boolean values (An empty string is same as `FALSE`).

## License

Copyright (c) 2016 Shogo Iwano
Licensed under the MIT license.
