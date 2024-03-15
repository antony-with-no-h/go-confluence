# go-confluence

Convert and publish Markdown to Confluence.

## Usage

```
Usage:
  go-confluence [command]

Available Commands:
  add         
  completion  Generate the autocompletion script for the specified shell
  edit        
  help        Help about any command

Flags:
  -f, --file string     Path to file containing Page Markdown
  -h, --help            help for go-confluence
  -p, --parent string   Title of the page that will act as the parent (e.g. Support, Backup and Restore)
  -s, --space string    The Confluence space where the page should be published (e.g. Engineering, QA)
  -t, --title string    Page title

Use "go-confluence [command] --help" for more information about a command.
```

### Examples

Publish a page titled *How-to Guides* in the Engineering space.

```
go-confluence add page \
    -s Engineering -t "How-to Guides" -f ./wiki/how_to_guides.md
```

Add a page nested under *How-to Guides* 

```
go-confluence add page \
    -s Engineering -p "How-to Guides" -t "Restore Database from Backup" \
    -f ./wiki/restore_database.md
```

Update the page

```
go-confluence edit page \
    -s Engineering -t "Restore Database from Backup" -f ./wiki/restore_database.md
```
