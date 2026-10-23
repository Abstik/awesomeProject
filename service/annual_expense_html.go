package service

import (
	"bytes"
	"golang.org/x/net/html"
	"net/url"
	"regexp"
	"strings"
)

var expenseStyle = regexp.MustCompile(`^(text-align|font-size|font-family|font-weight|font-style|text-decoration|color|background-color|line-height|text-indent|margin-left):[a-zA-Z0-9#.,%() "'-]+$`)
var expenseSpan = regexp.MustCompile("^[1-9][0-9]?$")
var expenseTags = map[string]bool{}

func init() {
	for _, tag := range strings.Fields("p br div span h1 h2 h3 h4 h5 h6 strong b em i u s del blockquote ul ol li table thead tbody tr td th hr a img sup sub") {
		expenseTags[tag] = true
	}
}
func safeExpenseURL(raw string, image bool) bool {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || strings.HasPrefix(strings.TrimSpace(raw), "//") {
		return false
	}
	if u.Scheme == "" {
		return strings.HasPrefix(raw, "/")
	}
	return u.Scheme == "https" || u.Scheme == "http" || (!image && u.Scheme == "mailto")
}
func sanitizeAnnualExpense(content string) string {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return ""
	}
	var clean func(*html.Node)
	clean = func(n *html.Node) {
		for child := n.FirstChild; child != nil; {
			next := child.NextSibling
			if child.Type == html.CommentNode {
				n.RemoveChild(child)
			} else if child.Type == html.ElementNode {
				if strings.Contains("|script|style|iframe|object|embed|svg|math|template|", "|"+child.Data+"|") {
					n.RemoveChild(child)
				} else {
					clean(child)
					if !expenseTags[child.Data] {
						for child.FirstChild != nil {
							c := child.FirstChild
							child.RemoveChild(c)
							n.InsertBefore(c, child)
						}
						n.RemoveChild(child)
					} else {
						attrs := []html.Attribute{}
						for _, a := range child.Attr {
							if a.Namespace != "" {
								continue
							}
							switch a.Key {
							case "style":
								styles := []string{}
								for _, s := range strings.Split(a.Val, ";") {
									s = strings.TrimSpace(s)
									if expenseStyle.MatchString(s) {
										styles = append(styles, s)
									}
								}
								if len(styles) > 0 {
									a.Val = strings.Join(styles, ";")
									attrs = append(attrs, a)
								}
							case "href":
								if child.Data == "a" && safeExpenseURL(a.Val, false) {
									attrs = append(attrs, a)
								}
							case "src":
								if child.Data == "img" && safeExpenseURL(a.Val, true) {
									attrs = append(attrs, a)
								}
							case "alt", "title":
								attrs = append(attrs, a)
							case "colspan", "rowspan":
								if expenseSpan.MatchString(a.Val) {
									attrs = append(attrs, a)
								}
							}
						}
						child.Attr = attrs
					}
				}
			}
			child = next
		}
	}
	clean(doc)
	var output bytes.Buffer
	html.Render(&output, doc)
	return output.String()
}
func annualExpenseText(content string) string {
	doc, _ := html.Parse(strings.NewReader(content))
	var out strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			out.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	if doc != nil {
		walk(doc)
	}
	return strings.ReplaceAll(out.String(), "\u200b", "")
}
