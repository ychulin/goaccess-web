package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	uploadDir = "./uploads"
	reportDir = "./reports"
)

// 主页 HTML
var tpl = `
<!DOCTYPE html>
<html>
<head>
    <title>GoAccess 日志分析</title>
</head>
<body>
    <h2>上传日志文件进行分析</h2>
    <form action="/upload" method="post" enctype="multipart/form-data">
        <input type="file" name="logfile" required>
        <input type="submit" value="上传并分析">
    </form>
    <h3>分析结果：</h3>
    <a href="/report.html" target="_blank">查看分析报告</a>
</body>
</html>
`

// 主页处理
func homeHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("index")
	t.Parse(tpl)
	t.Execute(w, nil)
}

// 文件上传处理
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 限制上传大小
	r.ParseMultipartForm(10 << 20) // 10MB
	file, handler, err := r.FormFile("logfile")
	if err != nil {
		http.Error(w, "无法获取文件", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 创建保存目录
	os.MkdirAll(uploadDir, os.ModePerm)
	os.MkdirAll(reportDir, os.ModePerm)

	// 目标文件路径
	filePath := filepath.Join(uploadDir, handler.Filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "无法保存文件", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()
	io.Copy(outFile, file)

	// 运行 GoAccess 进行日志分析
	reportPath := filepath.Join(reportDir, "report.html")
	cmd := exec.Command("goaccess", "-f", filePath, "--log-format=COMBINED", "-o", reportPath)
	err = cmd.Run()
	if err != nil {
		http.Error(w, "日志分析失败", http.StatusInternalServerError)
		return
	}

	// 跳转到报告页面
	http.Redirect(w, r, "/report.html", http.StatusFound)
}

// 提供 HTML 报告
func reportHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(reportDir, "report.html"))
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/report.html", reportHandler)

	fmt.Println("服务器启动，访问 http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
