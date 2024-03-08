# go-confluence

CLI to:

- Publish a file containing Confluence XHTML to a Confluence instance.
- Convert and publish a Markdown file to Confluence.

This project was created to strength my knowledge of Go, Go tooling: Cobra & Viper. And migrate some Python scripts to Go. The Markdown to XHTML conversion doesn't account for all available Confluence macros.

## Installation

If you have Go installed

```
go build
```

Create a file named `config` in `$XDG_CONFIG_HOME/go-confluence`. The path to search is setup using [UserConfigDir](https://pkg.go.dev/os#UserConfigDir).

Contents of the file should be:

```
AccessToken: Bearer 
Target: https://
CodeMacro:
  linenumbers: "true"
  theme: confluence
```

`AccessToken` is a [Confluence PAT](https://confluence.atlassian.com/enterprise/using-personal-access-tokens-1026032365.html) 

`CodeMacro` controls how all code macros will be themed and if to display line numbers.

## Usage

Commands currently implemented

- get
- post

### get [command]

Really a dummy function that I used to get going with Cobra. Can do some basic querying against the `/content` endpoint.

```
# unfiltered query
go-confluence get content

# add query strings
go-confluence get content spaceKey=engineering title="How-to guides"
```

### post [command]

From a file already written in Confluence XHTML

```
go-confluence add page --space <space> \
    --filename /home/user/wiki/template.txt \
    --title <title>
```

Convert a Markdown file and publish

```
go-confluence add page --space <space> \
    --filename /home/user/wiki/template.md \
    --title <title> --md
```
