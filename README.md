# Markdown internal & external links validation

## Overview

`MILV` is a bot that parses, checks and validates internal & external URLs links in markdown files. It can be used for [verification pull requests](#validate-pull-requests) and also as standalone library.

## Installation

```bash
$ go get -u -v github.com/kyma-incubator/milv
```

For the above command to work you must have [GoLang](https://golang.org/doc/install) installed

After installation, run the program with `milv` from anywhere in the file system.

### Run the code from source

If you want run the code without installation, run the following commands to get the source code, resolve external dependencies and build the project.

For this operations you must have also installed package manager [Dep](https://github.com/golang/dep).

```bash
git clone https://github.com/kyma-incubator/milv.git
cd milv
dep ensure
go build
```

## Usage

### Command line parameters

You can use the following parameters while using `milv` binary:

| Name                           | Description                                             | Default Value      |
| ------------------------------ | ------------------------------------------------------- | ------------------ |
| `-base-path`                   | Root directory of the repository                            | `""`               |
| `-config-file`                 | Configuration file for the bot. See the [**Config file**](#config-file) section for more details.  | `milv.config.yaml` |
| `-external-links-to-ignore`    | Comma-separated external links which will not be checked | `[]`               |
| `-internal-links-to-ignore`    | Comma-separated internal links which will not be checked | `[]`               |
| `-files-to-ignore`             | Comma-separated files which will not be checked          | `[]`               |
| `-allow-redirect`              | Redirects will be allowed                               | `false`            |
| `-request-repeats`             | Number of repeated request                               | `1`                |
| `-allow-code-blocks`           | Links in code blocks will be checked                          | `false`            |
| `-timeout`                     | Connection timeout (in seconds)                         | `30`               |
| `-ignore-external`             | External links to be ignored                                  | `false`            |
| `-ignore-internal`             | Internal links to be ignored                                   | `false`            |
| `-v`                           | Verbose logging                                  | `false`            |
| `-help` or `-h`                | Available parameters                               | n/a                |

Files to be checked are given as free parameters.

### Examples

- Checks all links, without matching `github.com` in external links, in `.md` files in current directory+subdirectories without files matching `vendor` in path:

```bash
milv -files-to-ignore="vendor" -external-links-to-ignore="github.com"
```

- Checks links only in `./README.md` and `./foo/bar.md` files:

```bash
milv ./README.md ./foo/bar.md
```

### Docker image

If you do not want to install `milv` and it's dependencies you can simply use Docker and Docker image:

```bash
docker run --rm -v $PWD:/milv:ro magicmatatjahu/milv:stability -base-path=/milv
```

## Config file

The configuration file allows for quick parameterization of the `milv` works. Config file must be a `.yaml` file.

Parameterization is very similar to using parameters in the `CLI`. However, you can configure files, located in subdirectories relative to the configuration file, separately with different config.

### Examples

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

your config file can look like this:

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

Before running the validation, `milv` removes the `./README.md` file from the files list. It then concatenates the values for `external-links-to-ignore` entries. For the `./src/foo.md` file mentioned in the example, the result will look as follows:

```yaml
external-links-to-ignore: ["localhost", "abc.com", "github.com"]
```

Similarly will be with `internal-links-to-ignore`.

If you have a config file and you use a `CLI`, then `milv` will automatically combine the parameters from file and consol.

#### Advanced configuration

> **NOTE**: For this example tree of project is the same as above.

Config file can look like this:

```yaml
external-links-to-ignore: ["localhost", "abc.com"]
internal-links-to-ignore: ["LICENSE"]
files-to-ignore: ["./README.md"]
request-repeats: 5
timeout: 45
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

In this example we can see that `milv` will globally check external links with 45 seconds timeout, also won't allow redirect and will allow checking links in code snippets and default times of request repeats is set 5.

`Milv` also allows to separately configurate files. Timeout in `./src/foo.md` file will be set to 30 seconds, links will be checking 3 times (if they will return error) and the links in code blocks won't be checked. However, a single link `https://github.com/kyma-incubator/milv` will be checking with 15 seconds timeout with the possibility of redirection.

## Troubleshooting links

The below table describes the types of errors during checking links and examples of how to solve them:

| Error                                                                                              | Solution example                                                                                                                                                                                                                                                              |
| -------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `404 Not Found`                                                                                    | Page doesn't exist - you have to change the external link to the correct one                                                                                                                                                                                                  |
| Error with formatting the link                                                                         | Correct the link. If the link contains variables or is used as an example, add it to the `external-links-to-ignore` or `internal-links-to-ignore` list.                                                                                                                                              |
| `The specified file doesn't exist`                                                                 | Change the relative path to the file to the correct one or use a absolute path (second solution is not recommended)                                                                                                                                                           |
| `The specified header doesn't exist in file`                                                       | Change the anchor link in `.md` file to the correct one. Sometimes `milv` give a hint (`Did you mean about <similar header>?`) of which header (existing in the file) is very similar to the given.                                                                           |
| `The specified anchor doesn't exist...` or `The specified anchor doesn't exist in website...`      | Check which anchors are on the external website and correct the specified anchor or remove the redirection to the given anchor. Sometimes `milv` give a hint (`Did you mean about <similar anchor>?`) of which anchor (existing in the website) is very similar to the given. |
| `Get <external link>: net/http: request canceled (Client.Timeout exceeded while awaiting headers)` | Increase net timeout to the all files, specific file or specific link or increase times of request repeats ([Here's](#advanced-configuration) how to do it)                                                                                                                   |
| `Get <external link>: EOF `                                                                        | Same as above or change the link to the other one (probably website doesn't exist)                                                                                                                                                                                            |
| Other types of errors and errors that contain the `no such host` or `timeout` words                   | It means that the website doesn't exist or you don't have access to it. You can change the link to another one, correct or remove it, or add it to the `external-links-to-ignore` or `internal-links-to-ignore` list.                                                               |

It is a good practice to add local or internal (in the local network) links to the global ignore list of external or internal links, such as `http://localhost`.

## Validate Pull Requests

`milv` can help you validate links in all `.md` files in whole repository when a pull request is created (or a commit is pushed).

### Jenkins

To use `milv` with Jenkins, connect your repo and create a [`Jenkinsfile`](https://jenkins.io/doc/book/pipeline/jenkinsfile/#creating-a-jenkinsfile) and add stage:

```groovy
stage("validate internal & external links") {
    workDir = pwd()
    sh "docker run --rm --dns=8.8.8.8 --dns=8.8.4.4 -v $workDir:/milv:ro magicmatatjahu/milv:0.0.6 -base-path=/milv"
}
```

## Other validators

In opensource community is available other links validation libraries written in JS, Ruby and others languages. Here are a few of note:

- [awesome_bot](https://github.com/dkhamsing/awesome_bot): validator written in Ruby. Allows for validation external and internal links in `.md` files.
- [remark-validate-links](https://github.com/remarkjs/remark-validate-links): validator written in JS. Allows for validation internal links in `.md` files.

## Contact

- [github.com/magicmatatjahu](https://github.com/magicmatatjahu)

## Contributing

If you want contribute this project, firstly read [CONTRIBUTING.md](CONTRIBUTING.md) file for details of submitting pull requests.

## License

This project is available under the MIT license. See the [LICENSE](LICENSE) file for more info.

## ToDo

- [ ] error handling
- [ ] refactor (new architecture)
- [ ] documentations
- [ ] possibility to validation remote repositories hosted on **GitHub**
- [ ] parse other type of files
- [x] add more commands like a: timeout for http.Get(), allow redirects or SSL
- [ ] landing page for project
