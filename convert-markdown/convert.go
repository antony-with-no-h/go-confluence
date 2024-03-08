package convert_markdown

import (
	"fmt"
	"io"
	"strings"

	"github.com/antony-with-no-h/go-confluence/config"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Macro struct {
	Theme       string
	Linenumbers bool
	Collapse    bool
}

func RenderHTML(b []byte, cfg config.Config) string {

	exts := parser.CommonExtensions | parser.HardLineBreak
	p := parser.NewWithExtensions(exts)

	doc := p.Parse(b)

	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: hook(cfg),
	}
	rend := html.NewRenderer(opts)

	//fmt.Printf("%s\n", markdown.Render(doc, rend))

	return fmt.Sprintf("%s", markdown.Render(doc, rend))
}

func hook(mTheme config.Config) html.RenderNodeFunc {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		//<ac:parameter ac:name=\"theme\">Confluence</ac:parameter>
		//<ac:parameter ac:name="title">bob</ac:parameter>
		//<ac:parameter ac:name="linenumbers">true</ac:parameter>

		if codeBlock, ok := node.(*ast.CodeBlock); ok {
			if entering {
				block := []string{
					"<ac:structured-macro ac:name=\"code\" ac:schema-version=\"1\">",
					fmt.Sprintf("<ac:parameter ac:name=\"language\">%s</ac:parameter>", theme(codeBlock)),
				}

				for k, v := range mTheme.CodeMacro {
					block = append(block,
						fmt.Sprintf("<ac:parameter ac:name=\"%s\">%s</ac:parameter>", k, v))
				}

				block = append(block,
					fmt.Sprintf("<ac:plain-text-body><![CDATA[%s]]></ac:plain-text-body></ac:structured-macro>", codeBlock.Literal))

				io.WriteString(w, strings.Join(block, ""))
			}

			return ast.GoToNext, true
		}
		return ast.GoToNext, false
	}
}

//func hook(theme *Macro) html.RenderNodeFunc {
//	return makeHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
//		if codeBlock, ok := node.(*ast.CodeBlock); ok {
//			confluenceCodeBlock(w, codeBlock, entering)
//			return ast.GoToNext, true
//		}
//		return ast.GoToNext, false
//	}
//}

func confluenceCodeBlock(w io.Writer, c *ast.CodeBlock, entering bool) {
	if entering {
		io.WriteString(w, fmt.Sprintf(`<ac:structured-macro ac:name="code" ac:schema-version="1">
<ac:parameter ac:name="language">%s</ac:parameter>
<ac:plain-text-body><![CDATA[%s]]></ac:plain-text-body></ac:structured-macro>`, theme(c), c.Literal))
	}
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
