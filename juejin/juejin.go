package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	markdown "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gocolly/colly"
)

var fileName = ""
var root = "./markdown"
var post = flag.String("post", "", "掘金文章ID")
var rootDir = flag.String("root", root, "文件保存的目录")

func main() {
	flag.Parse()
	c := colly.NewCollector(
		colly.Async(true),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnHTML(".article-title", func(e *colly.HTMLElement) {
		fileName = e.Text
		fileName = strings.TrimSpace(fileName)
	})

	c.OnHTML(".markdown-body", func(e *colly.HTMLElement) {
		reg := regexp.MustCompile(`data-`)
		html, _ := e.DOM.Html()
		htmlToMarkdown := convertHTMLToMarkdown(reg.ReplaceAllString(html, ""))
		writeFile(htmlToMarkdown)
	})

	_ = c.Visit("https://juejin.im/post/" + *post)
	c.Wait()
}

// 将Html转为Markdown
func convertHTMLToMarkdown(html string) string {
	converter := markdown.NewConverter("", true, nil)
	convertString, _ := converter.ConvertString(html)
	return convertString
}

// 写入文件
func writeFile(content string) {
	filePath := *rootDir + "/" + fileName + ".md"
	var file *os.File

	if !checkFileIsExist(root) {
		err := os.Mkdir(root, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	if checkFileIsExist(filePath) {
		// 如果文件存在，则删除
		err := os.Remove(filePath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 创建文件并写入内容
	fmt.Println("《" + fileName + "》" + " is downloaded on " + *rootDir + "/")
	file, _ = os.Create(filePath)
	n, _ := io.WriteString(file, "## "+fileName+"\n\n"+content)
	// 关闭文件
	defer file.Close()
	if n == 0 {
		return
	}
}

// 检查文件是否存在
func checkFileIsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}
