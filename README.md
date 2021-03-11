# MILV

## Overview

MILV stands for "Markdown internal and external links validation." It is a tool that parses, checks, and validates internal and external URL links in Markdown files.
You can use it either for verifying pull requests or as a standalone library.

## Prerequisites

To use MILV, you must have [GoLang](https://golang.org/doc/install) in version 1.15 or higher installed.

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

Build a Docker image out of MILV to use it in a continuous integration pipeline.

```bash
make build-image
```

You can use the Docker image in the same way as the MILV binary. To do that, mount a directory that MILV should validate.

```
cd ..
docker run -v $PWD:/milv milv -base-path milv
```

## Usage

### Definitions

MILV's logic is based on this basic distinction:
- **Internal link** is the link to the local resource, header, or any other file.
- **External link** is the link to the HTTP resource.

### Command line parameters

You can use the following parameters when using the MILV binary:

| Name                           | Description                                                 | Default Value      |
| ------------------------------ | ------------------------------------------------------------| ------------------ |
| `-base-path`                   | Root directory of the repository                            | `""`               |
| `-backoff`                     | Backoff timeout                                             | `"1s"`             |
| `-config-file`                 | Configuration file for the bot. See the [**Configuration file**](/docs/configuration-file.md) for more details.  | `milv.config.yaml` |
| `-external-links-to-ignore`    | Comma-separated external links which MILV must not check    | `[]`               |
| `-internal-links-to-ignore`    | Comma-separated internal links which MILV must not check    | `[]`               |
| `-files-to-ignore`             | Comma-separated files which MILV must not check            | `[]`               |
| `-allow-redirect`              | Redirects should be allowed                                   | `false`            |
| `-request-repeats`             | Number of repeated request                                  | `1`                |
| `-allow-code-blocks`           | Validating links in code blocks should be allowed                        | `false`            |
| `-timeout`                     | Connection timeout (in seconds)                             | `30`               |
| `-ignore-external`             | External links that MILV must ignore                                | `false`            |
| `-ignore-internal`             | Internal links that MILV must ignore                                 | `false`            |
| `-v`                           | Verbose logging                                             | `false`            |
| `-help` or `-h`                | Available parameters                                        |  n/a                |

Files to be checked are given as free parameters.

See these examples:

- Use this command to check all links except for the `vendor` directory in the project and external links containing the `github.com` address.

  ```bash
  milv -files-to-ignore="vendor" -external-links-to-ignore="github.com"
  ```

- Use this command to check all links in the `./README.md` and `./foo/bar.md` files:

  ```bash
  milv ./README.md ./foo/bar.md
  ```

### Configuration File

MILV relies on the `milv.config.yaml` configuration file in which you define rules and exceptions for MILV, stating which files and types of links it should validate or ignore. See the [**Configuration file**](/docs/configuration-file.md) for sample `milv.config.yaml` files and and a list of parameters you can use to configure such a file.

### Typical errors

The table describes types of errors MILV can return while checking the links and sample solutions to these issues:

| Error                                                                                              | Solution example                                                                                                                                                                                                                                                              |
| --- | --- |
| `404 Not Found`                                                                                    | This page doesn't exist. Change the external link to the correct one.                        |
| Error with link formatting                                                                    | Correct the link. If the link contains variables or is used as an example, add it to the **external-links-to-ignore** or **internal-links-to-ignore** list.  |
| `The specified file doesn't exist`                                                                 | Change the relative path to the file to the correct one. Alternatively, use an absolute path. |
| `The specified header doesn't exist in this file`                                                       | Change the anchor link in the MD file to the correct one. MILV sometimes gives a hint (`Did you mean {similar header}?`) and points to an existing header in the file that is very similar to the one provided.    |
| `The specified anchor doesn't exist` or `The specified anchor doesn't exist on the website`      | Check which anchors are on the external website and correct the specified anchor or remove the redirection to the given anchor. MILV sometimes gives a hint (`Did you mean {similar anchor}?`) and points to an existing header in the file that is very similar to the one provided. |
| `Get {external link}: net/http: request canceled (Client.Timeout exceeded while awaiting headers)` | Increase net timeout for all files, a specific file, or a specific link. Alternatively, increase the the value for **request-repeats**. See the [**Configuration file**](/docs/configuration-file.md) for more details.  |
| `Get {external link}: EOF`                                                                        | Follow the already mentioned steps. You can also change the link to another one as it is possible that the website doesn't exist. |
| Other types of errors, such as errors that contain the `no such host` or `timeout` words                | It means that the website doesn't exist or you don't have access to it. You can change the link to another one, correct, or remove it. Alternatively, add the link to the **external-links-to-ignore** or **internal-links-to-ignore** list.   |

It is considered a good practice to add external local links (in the local network) to the global ignore list of external links, such as `http://localhost`.

## Development

If you want to contribute to this project, read the [`CONTRIBUTING.md`](CONTRIBUTING.md) file for hints on how to submit pull requests.
