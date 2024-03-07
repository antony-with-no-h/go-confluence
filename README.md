# go-confluence

Some basic CLI tooling for the Confluence "Data Center REST API".

## Usage

### Create page

Use a file already written in Confluence XHTML

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

## Project Goals

- Strengthen existing go knowledge
- Use Cobra and Viper for the first time
- Migrate some existing Python scripts to Go

