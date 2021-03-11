# Configuration File

The configuration file is where you parameterize MILV, and specify which files or links MILV should ignore in a given project.
A configuration file must be a YAML file and should be named `milv.config.yaml`.

You can overwrite the name of the configuration file with the `-config-file` command line parameter.
You can use a list of such [command line arguments](../README.md#command-line-parameters) to provide additional configuration for MILV or to overwrite the configuration file parameters.

Place the configuration file at the root of your repository. See a sample project file structure:

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

## Configurable parameters

`milv.config.yaml` can take the following parameters:

| Parameter                           | Description                                                | Type | Default Value      |
| ------------------------------ | ------------------------------------------------------------| ------|------------ |
| **backoff**| Amount of time MILV must wait for the next external link validation when the server responds with the `429` status code (`Too many requests`) | duration | `1s` |
| **external-links-to-ignore** | List of external links for MILV to ignore | array of strings | n/a |
| **internal-links-to-ignore** | List of internal links for MILV to ignore | array of strings| n/a |
| **files-to-ignore** | List of files and directories in which MILV won't check any links | array of strings | n/a |
| **files-to-ignore-internal-links-in** | List of files and directories in which MILV won't check internal links | array of strings | n/a |
| **timeout** | Timeout for the HTTP external links check | integer | `30` |
| **request-repeats** | Number of HTTP tries when validating external links | integer | `1` |
| **allow-redirect** | Parameter specifying if MILV should follow redirects in the whole project | boolean  | `false` |
| **allow-code-blocks** | Parameter specifying if MILV should check links in code blocks |  boolean | `false` |
| **ignore-external** | External links will be ignored | boolean | `false` |
| **ignore-internal** | Internal links will be ignored | boolean | `false` |
| **files** | List of files for which MILV must apply different settings | n/a |
| **files.path** | Path to the file | string | n/a |
| **files.links** | List of link settings for the file | array of objects | n/a |
| **files.links.path** | Link name | string | n/a |
| **files.links.config** | Configuration of a specific link in the file | object | n/a |
| **files.links.config.timeout** | Timeout for the HTTP external links check | integer | `30` |
| **files.links.config.request-repeats** | Number of HTTP tries when validating external links | integer | `1` |
| **files.links.config.allow-redirect** | Parameter specifying if MILV should follow redirects for the given link | boolean | `false` |
| **files.config** | Configuration of a specific file | object | n/a |
| **files.config.backoff** | Amount of time MILV must wait for the next external link validation when the server responds with the `429` status code (`Too many requests`) | duration | `1s` |
| **files.config.external-links-to-ignore** | Specific external links for MILV to ignore | array of strings | n/a |
| **files.config.internal-links-to-ignore** | Specific internal links for MILV to ignore | array of strings | n/a |
| **files.config.timeout** | Timeout for the HTTP external links check | integer | `30` |
| **files.config.request-repeats** | Number of HTTP tries when validating external links | integer | `1` |
| **files.config.allow-redirect** | Parameter specifying if MILV should follow redirects for links in this file | boolean | `false` |
| **files.config.allow-code-blocks** | Parameter specifying if MILV should check links in code blocks in this file | boolean | `false` |
| **files.config.ignore-external** | MILV will ignore all external links in this file | boolean | `false` |
| **files.config.ignore-internal** | MILV will ignore all internal links in this file | boolean | `false` |

## Basic configuration file

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

## Advanced configuration file

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
- Waits 2 seconds before making another call if the server responds with the `429` status code (`Too many requests`).
- Follows redirects.
- Checks links in code snippets.
- Makes a maximum of 5 requests in case of an error.

For `/src/foo.md`, MILV:
- Timeouts after 30 seconds.
- Makes a maximum of 3 requests in case of an error.
- Ignores links in code blocks.
- For the `https://github.com/kyma-incubator/milv` link, MILV will timeout after 15 seconds and follow the redirects.
