# Markdown internal & external links validation

## Overview

MILV is a program that parses, checks, and validates internal and external URL links in Markdown files.
You can use it either for verifying pull requests or as a standalone library.

## Prerequisites

You must have [GoLang](https://golang.org/doc/install) in version 1.15 or higher installed.

## Installation

### Build a binary

Run the following commands to get the source code, resolve external dependencies, and build the project.
The output binary named `milv` will be created in the project directory.

```bash
git clone https://github.com/kyma-incubator/milv.git
cd milv
make build
```

### Build a Docker image

You can build a Docker image out of MILV to use it in a continuous integration pipeline.
```bash
make build-image
```

You can use it as normal `milv` binary, you only need to mount volume, example:
```
cd ..
docker run -v $PWD:/milv milv -base-path milv
```

### Definitions

MILV's logic is based on this basic distinction:
- **Internal link** is the link to the local resource, header, or any other file.
- **External link** is the link to the HTTP resource.

### Command line parameters

You can use the following parameters while using `milv` binary:

| Name                           | Description                                                 | Default Value      |
| ------------------------------ | ------------------------------------------------------------| ------------------ |
| `-base-path`                   | Root directory of the repository                            | `""`               |
| `-backoff`                     | Backoff timeout                                             | `"1s"`             |
| `-config-file`                 | Configuration file for the bot. See the [**Config file**](#configuration-file) section for more details.  | `milv.config.yaml` |
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

The configuration file allows you to parameterize MILV, and specify clearly which files or links MILV should ignore.
A configuration file must be a `.yaml` file and should be named `milv.config.yaml`, unless you overwrite it with the `-config-file` command line parameter.

You can use the configuration file with command line arguments. 
You can use the command line parameters to provide additional configuration for MILV or to overwrite the configuration file parameters.

#### Specification of configuration file

| Parameter                           | Description                                                | Type | Default Value      |
| ------------------------------ | ------------------------------------------------------------| ------|------------ |
| backoff| Amount of time MILV must wait for the next external link validation when the server responds with the 429 status code (`Too many requests`) | duration | 1s |
| external-links-to-ignore | List of external links which have to be ignored | array of strings | |
| internal-links-to-ignore | List of internal links for MILV to ignore | array of strings| |
| files-to-ignore | List of files and directories in which MILV won't check any links | array of strings | |
| files-to-ignore-internal-links-in | List of files and directories in which MILV won't check internal links | array of strings | |
| timeout | Timeout for the HTTP external links check | integer | 30 |
| request-repeats | Number of HTTP tries when validating external links | integer | 1 |
| allow-redirect | Allow following redirects | boolean  | false |
| allow-code-blocks | Allow MILV to check links in code blocks |  boolean | false |
| ignore-external | External links will be ignored | boolean | false |
| ignore-internal | Internal links will be ignored | boolean | false |
| files | List of files which should be treated with different settings | array of objects | |
| files.path | Path to the file | string | |
| files.links | List of link settings for the file | array of objects | | 
| files.links.path | Link name | string | |
| files.links.config | Configuration of a specific link in the file | object | |
| files.links.config.timeout | Timeout for the HTTP external links check | integer | 30 |
| files.links.config.request-repeats | Number of HTTP tries when validating external links | integer | 1 |
| files.links.config.allow-redirect | Allow following redirects | boolean | false |
| files.config | Configuration of a specific file | object | |
| files.config.backoff | Amount of time MILV must wait for the next external link validation when the server responds with the 429 status code (`Too many requests`) | duration | 1s | 

| files.config.external-links-to-ignore | Files in which MILV must ignore all external links | array of strings | |

| files.config.internal-links-to-ignore | Files in which MILV must ignore all internal links | array of strings | | 
| files.config.timeout | Timeout for the HTTP external links check | integer | 30 |
| files.config.request-repeats | Number of HTTP tries when validating external links | integer | 1 |
| files.config.allow-redirect | Allow following redirects | boolean | false |
| files.config.allow-code-blocks | Links in code blocks will be checked in the file | boolean | false |
| files.config.ignore-external | External links will be ignored in the file | boolean | false |
| files.config.ignore-internal | Internal links will be ignored in the file | boolean | false |


## Usage

### Command line parameters

Use this command to check all links except for the `vendor` directory in the project and external links containing the `github.com` address.

```bash
milv -files-to-ignore="vendor" -external-links-to-ignore="github.com"
```

Use this command to check all links in the `./README.md` and `./foo/bar.md` files:

```bash
milv ./README.md ./foo/bar.md
```

### Basic configuration file

Your project file structure should look as follows:

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

See a sample configuration file:

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

Before running the validation, MILV removes the `./README.md` file from the files list.

For the `./src/foo.md` file mentioned in the example, MILV concatenates the values for **external-links-to-ignore** entries.
The list of ignored external links will look as follows:
```yaml
external-links-to-ignore: ["localhost", "abc.com", "github.com"]
```

The same mechanism applies to the **internal-links-to-ignore** parameter.

### Advanced configuration file

> **NOTE**: For this example tree of project is the same as above.

See a sample configuration file:

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

Having this configuration, MILV globally:
- Checks external links with the 45 seconds timeout.
-  Waits 2 seconds before making another call if the server responds with the 429 status code (`Too many requests`).
- won't follow redirects
- Checks links in code snippets.
- Makes a maximum of 5 requests in case of an error.

For `/src/foo.md`, MILV:
- Timeouts after 30 seconds.
- Makes a maximum of 3 requests in case of an error.
- Ignores links in code blocks.
- For the `https://github.com/kyma-incubator/milv` link, MILV will timeout after 15 seconds and follow redirections.

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
