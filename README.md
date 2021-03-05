# Markdown internal & external links validation

## Overview

`MILV` is a program that parses, checks and validates internal & external URLs links in markdown files.
It can be used for verification pull requests and as standalone library.

## Prerequisites

You have to have [GoLang](https://golang.org/doc/install) in version 1.15 or above.) installed

## Installation

### Build binary

Run the following commands to get the source code, resolve external dependencies and build the project.
The output binary will be in project directory with name: `milv`.

```bash
git clone https://github.com/kyma-incubator/milv.git
cd milv
make build
```

### Build docker image

There is an option to build docker image.
```bash
make build-image
```

You can use it as normal `milv` binary, you only need to mount volume, example:
```
cd ..
docker run -v $PWD:/milv milv -base-path milv
```

### Definitions

Internal link - is the link to the local resource, header, other file.
External link - is the link to the HTTP resource

### Command line parameters

You can use the following parameters while using `milv` binary:

| Name                           | Description                                                 | Default Value      |
| ------------------------------ | ------------------------------------------------------------| ------------------ |
| `-base-path`                   | Root directory of the repository                            | `""`               |
| `-backoff`                     | Backoff timeout                                             | `"1s"`             |
| `-config-file`                 | Configuration file for the bot. See the [**Config file**](#config-file) section for more details.  | `milv.config.yaml` |
| `-external-links-to-ignore`    | Comma-separated external links which will not be checked    | `[]`               |
| `-internal-links-to-ignore`    | Comma-separated internal links which will not be checked    | `[]`               |
| `-files-to-ignore`             | Comma-separated files which will not be checked             | `[]`               |
| `-allow-redirect`              | Redirects will be allowed                                   | `false`            |
| `-request-repeats`             | Number of repeated request                                  | `1`                |
| `-allow-code-blocks`           | Links in code blocks will be checked                        | `false`            |
| `-timeout`                     | Connection timeout (in seconds)                             | `30`               |
| `-ignore-external`             | External links to be ignored                                | `false`            |
| `-ignore-internal`             | Internal links to be ignored                                | `false`            |
| `-v`                           | Verbose logging                                             | `false`            |
| `-help` or `-h`                | Available parameters                                        | n/a                |

Files to be checked are given as free parameters.

### Configuration file

The configuration file allows for more accurate parameterization of `milv` by having config per file.
Config file must be a `.yaml` file and should be named `milv.config.yaml` if not overwritten by `-config-file` command line parameter.
Full documentation of all fields is in [json.schema](./json.schema).

You can use the configuration file with command line arguments. Those settings will be combined.
The command line parameters overwrites the configuration file parameters.

### Usage

#### Command line parameters

- Checks all links, without matching `github.com` in all external links and without checking `vendor` directory in project.

```bash
milv -files-to-ignore="vendor" -external-links-to-ignore="github.com"
```

- Checks all links only in `./README.md` and `./foo/bar.md` files:

```bash
milv ./README.md ./foo/bar.md
```

#### Configuration file

##### Basic example
If your tree of your project look like this:

```
├── README.md
├── LICENSE
├── main.go
├── milv.config.yaml
└── src
    ├── file.go
    ├── file_test.go
    ├── foo.md
    └── some_dir
              └── bar.md
```

given config file:

```yaml
external-links-to-ignore: ["localhost", "abc.com"]
internal-links-to-ignore: ["LICENSE"]
files-to-ignore: ["./README.md"]
files:
  - path: "./src/foo.md"
    config:
      external-links-to-ignore: ["github.com"]
      internal-links-to-ignore: ["#contributing"]
```

Before running the validation, `milv` removes the `./README.md` file from the files list.

For the `./src/foo.md` file mentioned in the example, milv concatenates the values for `external-links-to-ignore` entries.
The list of ignored external links will be:
```yaml
external-links-to-ignore: ["localhost", "abc.com", "github.com"]
```

The same mechanism applies to `internal-links-to-ignore` parameter.

#### Advanced configuration

> **NOTE**: For this example tree of project is the same as above.

given config file:

```yaml
external-links-to-ignore: ["localhost", "abc.com"]
internal-links-to-ignore: ["LICENSE"]
files-to-ignore: ["./README.md"]
request-repeats: 5
timeout: 45
backoff: 2s
allow-redirect: false
allow-code-blocks: true
files:
  - path: "./src/foo.md"
    config:
      external-links-to-ignore: ["google.com"]
      internal-links-to-ignore: ["#contributing"]
      request-repeats: 3
      timeout: 30
      allow-code-blocks: false
    links:
      - path: "https://github.com/kyma-incubator/milv"
        config:
          timeout: 15
          allow-redirect: true
```

Having this configuration, the `milv` globally:
- will check external links with 45 seconds timeout
- will wait 2 second before doing another call if server responds with 429 status code (Too many requests)
- won't follow redirects
- will allow checking links in code snippets
- will do maximum of 5 requests in case of an error.

For `/src/foo.md` the `milv`:
- will timeout after 30 seconds
- will do maximum of 3 requests in case of an error.
- will ignore links in code blocks
- for link `https://github.com/kyma-incubator/milv` will timeout after 15 seconds and will follow redirections.

### Typical errors

The below table describes the types of errors during checking links and examples of how to solve them:

| Error                                                                                              | Solution example                                                                                                                                                                                                                                                              |
| -------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `404 Not Found`                                                                                    | Page doesn't exist - you have to change the external link to the correct one                                                                                                                                                                                                  |
| Error with formatting the link                                                                     | Correct the link. If the link contains variables or is used as an example, add it to the `external-links-to-ignore` or `internal-links-to-ignore` list.                                                                                                                                              |
| `The specified file doesn't exist`                                                                 | Change the relative path to the file to the correct one or use a absolute path (second solution is not recommended)                                                                                                                                                           |
| `The specified header doesn't exist in file`                                                       | Change the anchor link in `.md` file to the correct one. Sometimes `milv` give a hint (`Did you mean about <similar header>?`) of which header (existing in the file) is very similar to the given.                                                                           |
| `The specified anchor doesn't exist...` or `The specified anchor doesn't exist in website...`      | Check which anchors are on the external website and correct the specified anchor or remove the redirection to the given anchor. Sometimes `milv` give a hint (`Did you mean about <similar anchor>?`) of which anchor (existing in the website) is very similar to the given. |
| `Get <external link>: net/http: request canceled (Client.Timeout exceeded while awaiting headers)` | Increase net timeout to the all files, specific file or specific link or increase times of request repeats ([Here's](#advanced-configuration) how to do it)                                                                                                                   |
| `Get <external link>: EOF `                                                                        | Same as above or change the link to the other one (probably website doesn't exist)                                                                                                                                                                                            |
| Other types of errors and errors that contain the `no such host` or `timeout` words                | It means that the website doesn't exist or you don't have access to it. You can change the link to another one, correct or remove it, or add it to the `external-links-to-ignore` or `internal-links-to-ignore` list.                                                               |

It is a good practice to add external local links (in the local network) to the global ignore list of external links, such as `http://localhost`.

## Development

If you want contribute this project, firstly read [CONTRIBUTING.md](CONTRIBUTING.md) file for details of submitting pull requests.

## License

This project is available under the MIT license. See the [LICENSE](LICENSE) file for more info.
