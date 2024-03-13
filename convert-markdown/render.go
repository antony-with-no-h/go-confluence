package convert_markdown

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/antony-with-no-h/go-confluence/config"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func RenderHTML(b []byte, cfg config.Config) (string, []byte) {
	exts := parser.CommonExtensions
	p := parser.NewWithExtensions(exts)
	p.Opts.ParserHook = parserNoteHook

	doc := p.Parse(b)

	htmlFlags := html.UseXHTML
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: hook(cfg),
	}
	rend := html.NewRenderer(opts)

	return fmt.Sprintf("%s", markdown.Render(doc, rend)), markdown.Render(doc, rend)
}

func hook(cfg config.Config) html.RenderNodeFunc {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		if codeBlock, ok := node.(*ast.CodeBlock); ok {
			if entering {
				code(cfg, w, codeBlock)
			}

			return ast.GoToNext, true
		}

		if noteBlock, ok := node.(*Note); ok {
			if entering {
				note(cfg, w, noteBlock)
			}

			return ast.GoToNext, true
		}
		return ast.GoToNext, false
	}
}

func code(cfg config.Config, w io.Writer, node *ast.CodeBlock) {
	block := []string{
		"<ac:structured-macro ac:name=\"code\" ac:schema-version=\"1\">",
		fmt.Sprintf("<ac:parameter ac:name=\"language\">%s</ac:parameter>", theme(node)),
	}

	for k, v := range cfg.CodeMacro {
		block = append(block,
			fmt.Sprintf("<ac:parameter ac:name=\"%s\">%s</ac:parameter>", k, v))
	}

	block = append(block,
		fmt.Sprintf("<ac:plain-text-body><![CDATA[%s]]></ac:plain-text-body></ac:structured-macro>", node.Literal))

	io.WriteString(w, strings.Join(block, ""))
}

func note(cfg config.Config, w io.Writer, node *Note) {
	var body []string
	re := regexp.MustCompile(`^\s{3,}`)

	//for _, line := range bytes.SplitAfterN(node.Notes, []byte("\n"), 1) {
	for _, line := range bytes.Split(node.Notes, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		body = append(body, string(re.ReplaceAll(line, []byte(""))))
	}

	io.WriteString(w, fmt.Sprintf(`<ac:structured-macro ac:name="%s" ac:schema-version="1">                            
<ac:rich-text-body>%s</ac:rich-text-body></ac:structured-macro>`, string(node.Type), strings.Join(body, "\n")))
}

func theme(c *ast.CodeBlock) string {
	// aliases from github.com/github-linguist/linguist/blob/master/lib/linguist/languages.yml
	syntaxMap := map[string]string{
		"actionscript 3":  "actionscript3",
		"actionscript3":   "actionscript3",
		"as3":             "actionscript3",
		"applescript":     "applescript",
		"scpt":            "applescript",
		"sh":              "bash",
		"shell-script":    "bash",
		"bash":            "bash",
		"zsh":             "bash",
		"csharp":          "c#",
		"cake":            "c#",
		"cakescript":      "c#",
		"cpp":             "cpp",
		"css":             "css",
		"cfm":             "coldfusion",
		"cfml":            "coldfusion",
		"coldfusion html": "coldfusion",
		"delphi":          "delphi",
		"objectpascal":    "delphi",
		"diff":            "diff",
		"erlang":          "erl",
		"erl":             "erl",
		"groovy":          "groovy",
		"html":            "xml",
		"xhtml":           "xml",
		"xml":             "xml",
		"rss":             "xml",
		"xsd":             "xml",
		"wsdl":            "xml",
		"java":            "java",
		"js":              "js",
		"node":            "js",
		"php":             "php",
		"perl":            "perl",
		"none":            "text",
		"fundamental":     "text",
		"plain text":      "text",
		"powershell":      "powershell",
		"posh":            "powershell",
		"pwsh":            "powershell",
		"python":          "py",
		"jruby":           "ruby",
		"macruby":         "ruby",
		"rake":            "ruby",
		"rb":              "ruby",
		"rbx":             "ruby",
		"sql":             "sql",
		"sass":            "sass",
		"scala":           "scala",
		"visual basic":    "vb",
		"vbnet":           "vb",
		"vb .net":         "vb",
		"vb.net":          "vb",
		"yaml":            "yml",
		"yml":             "yml",
		"json":            "json",
		"geojson":         "json",
		"jsonl":           "json",
		"topojson":        "json",
	}

	syntax := fmt.Sprintf("%s", c.Info)
	if lang, ok := syntaxMap[syntax]; ok {
		return lang
	}

	return "text"
}
