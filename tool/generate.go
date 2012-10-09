package main

import (
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "github.com/russross/blackfriday"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strings"
)

var cacheDir = "/tmp/gobyexample-cache"
var siteDir = "site"

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func filterStrings(vs []string, f func(string) bool) []string {
    vsf := make([]string, 0)
    for _, v := range vs {
        if f(v) {
            vsf = append(vsf, v)
        }
    }
    return vsf
}

func ensureDir(dir string) {
    err := os.MkdirAll(dir, 0700)
    check(err)
}

func pipe(bin string, arg []string, src string) []byte {
    cmd := exec.Command(bin, arg...)
    in, _ := cmd.StdinPipe()
    out, _ := cmd.StdoutPipe()
    cmd.Start()
    in.Write([]byte(src))
    in.Close()
    bytes, _ := ioutil.ReadAll(out)
    err := cmd.Wait()
    check(err)
    return bytes
}

func sha1Sum(s string) string {
    h := sha1.New()
    h.Write([]byte(s))
    b := h.Sum(nil)
    return hex.EncodeToString(b)
}

func mustReadFile(path string) string {
    bytes, err := ioutil.ReadFile(path)
    check(err)
    return string(bytes)
}

func cachedPygmentize(lex string, src string) string {
    ensureDir(cacheDir)
    arg := []string{"-l", lex, "-f", "html"}
    bin := "/usr/local/bin/pygmentize"
    cachePath := cacheDir + "/pygmentize-" + strings.Join(arg, "-") + "-" + sha1Sum(src)
    cacheBytes, cacheErr := ioutil.ReadFile(cachePath)
    if cacheErr == nil {
        return string(cacheBytes)
    }
    renderBytes := pipe(bin, arg, src)
    writeErr := ioutil.WriteFile(cachePath, renderBytes, 0600)
    check(writeErr)
    return string(renderBytes)
}

func markdown(src string) string {
    return string(blackfriday.MarkdownCommon([]byte(src)))
}

func readLines(path string) []string {
    src := mustReadFile(path)
    return strings.Split(src, "\n")
}

func mustGlob(glob string) []string {
    paths, err := filepath.Glob(glob)
    check(err)
    return paths
}

func whichLexer(path string) string {
    if strings.HasSuffix(path, ".go") {
        return "go"
    } else if strings.HasSuffix(path, ".sh") {
        return "console"
    }
    panic("No lexer for " + path)
    return ""
}

func debug(msg string) {
    if os.Getenv("DEBUG") == "1" {
        fmt.Fprintln(os.Stderr, msg)
    }
}

var docsPat = regexp.MustCompile("^\\s*(\\/\\/|#)\\s")
var headerPat = regexp.MustCompile("^\\s*(\\/\\/|#)\\s#+\\s")
var todoPat = regexp.MustCompile("\\/\\/ todo: ")

type seg struct {
    docs, code, docsRendered, codeRendered string
}

func parseSegs(sourcePath string) []*seg {
    lines := readLines(sourcePath)
    segs := []*seg{}
    segs = append(segs, &seg{code: "", docs: ""})
    lastSeen := ""
    for _, line := range lines {
        if todoPat.MatchString(line) {
            continue
        }
        headerMatch := headerPat.MatchString(line)
        docsMatch := docsPat.MatchString(line)
        emptyMatch := line == ""
        lastSeg := segs[len(segs)-1]
        lastHeader := lastSeen == "header"
        lastDocs := lastSeen == "docs"
        newHeader := lastSeen != "header" && lastSeg.docs != ""
        newDocs := lastSeen == "code" || lastSeen == "header"
        newCode := (lastSeen != "code" && lastSeg.code != "") || lastSeen == "header"
        if newHeader || newDocs || newCode {
            debug("NEWSEG")
        }
        if headerMatch || (emptyMatch && lastHeader) {
            trimmed := docsPat.ReplaceAllString(line, "")
            if newHeader {
                newSeg := seg{docs: trimmed, code: ""}
                segs = append(segs, &newSeg)
            } else {
                lastSeg.docs = lastSeg.docs + "\n" + trimmed
            }
            debug("HEAD")
            lastSeen = "header"
        } else if docsMatch || (emptyMatch && lastDocs) {
            trimmed := docsPat.ReplaceAllString(line, "")
            if newDocs {
                debug("NEWSEG")
                newSeg := seg{docs: trimmed, code: ""}
                segs = append(segs, &newSeg)
            } else {
                lastSeg.docs = lastSeg.docs + "\n" + trimmed
            }
            debug("DOCS")
            lastSeen = "docs"
        } else {
            if newCode {
                newSeg := seg{docs: "", code: line}
                segs = append(segs, &newSeg)
            } else {
                lastSeg.code = lastSeg.code + "\n" + line
            }
            debug("CODE")
            lastSeen = "code"
        }
    }
    return append(segs, &seg{code: "", docs: ""})
}

func parseAndRenderSegs(sourcePath string) []*seg {
    segs := parseSegs(sourcePath)
    lexer := whichLexer(sourcePath)
    for _, seg := range segs {
        if seg.docs != "" {
            seg.docsRendered = markdown(seg.docs)
        }
        if seg.code != "" {
            seg.codeRendered = cachedPygmentize(lexer, seg.code)
        }
    }
    return segs
}

func generateIndex() {
    indexF, err := os.Create(siteDir + "/index.html")
    check(err)
    fmt.Fprint(indexF,
        `<!DOCTYPE html>
         <html>
           <head>
             <meta http-eqiv="content-type" content="text/html;charset=utf-8">
             <title>Go by Example</title>
             <link rel=stylesheet href="../style/site.css">
           </head>
           <body>
           <div class="chapter" id="contents"><h2>Contents</h2><ul>`)
    chapterIds := readLines("meta/contents.txt")
    for _, chapterId := range chapterIds {
        chapterLines := readLines("src/" + chapterId + "/" + chapterId + ".go")
        chapterName := chapterLines[0][6:]
        fmt.Fprintf(indexF, `<li><a href="%s.html">%s</a></li>`, chapterId, chapterName)
    }
    fmt.Fprint(indexF, `</ul></div></body></html>`)
}

func generateChapters() {
    chapterIds := readLines("meta/contents.txt")
    for _, chapterId := range chapterIds {
        chapterLines := readLines("src/" + chapterId + "/" + chapterId + ".go")
        chapterName := chapterLines[0][6:]
        chapterPath := "src/" + chapterId
        sourcePaths := mustGlob(chapterPath + "/*")
        chapterF, err := os.Create(siteDir + "/" + chapterId + ".html")
        check(err)
        fmt.Fprintf(chapterF,
            `<!DOCTYPE html>
             <html>
               <head>
                 <meta http-eqiv="content-type" content="text/html;charset=utf-8">
                 <title>Go by Example: %s</title>
                 <link rel=stylesheet href="../style/site.css">
               </head>
               <body>
                 <div class="chapter" id="%s">
                   <table cellspacing="0" cellpadding="0"><tbody>`,
            chapterName, chapterId)
        for _, sourcePath := range sourcePaths {
            segs := parseAndRenderSegs(sourcePath)
            for _, seg := range segs {
                codeClasses := "code"
                if seg.code == "" {
                    codeClasses = codeClasses + " empty"
                }
                fmt.Fprintf(chapterF,
                    `<tr>
                     <td class=docs>%s</td>
                     <td class="%s">%s</td>
                     </tr>`,
                    seg.docsRendered, codeClasses, seg.codeRendered)
            }
        }
        fmt.Fprint(chapterF, `</tbody></table></div></body></html>`)
    }
}

func main() {
    ensureDir(siteDir)
    generateIndex()
    generateChapters()
}
