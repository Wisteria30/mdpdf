package main

import (
    "fmt"
    "os/exec"
    "log"
    "strings"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// 拡張子チェック
func extCheck(checks []string, filename string) bool {
    for _, c := range checks {
        if filepath.Ext(filename) == c {
            return true
        }
    }
    return false
}


func main() {
    router := gin.Default()
    const PATH = "resource/"
    
    router.LoadHTMLGlob("public/*.html")
    
    router.GET("/", func(c *gin.Context) {
        c.HTML(200, "index.html", gin.H{})
    })

	router.POST("/upload", func(c *gin.Context) {
        convertType := c.PostForm("type")

        // imageファイル
		form, m_err := c.MultipartForm()
		if m_err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", m_err.Error()))
			return
		}
		image_files := form.File["image_files"]

		for _, image_file := range image_files {
            image_filename := filepath.Base(image_file.Filename)
            
            if !extCheck([]string{".png", "jpg", "jpeg"}, image_filename) {
                c.String(http.StatusBadRequest, fmt.Sprintf("upload file is not Image file."))
			    return
            }
			if err := c.SaveUploadedFile(image_file, PATH + image_filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload image file err: %s", err.Error()))
				return
			}
		}

        // mdファイル
		md_file, f_err := c.FormFile("md_file")
		if f_err != nil {
            c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", f_err.Error()))
			return
		}
        md_filename := filepath.Base(md_file.Filename)
        // md以外をはじく
        if !extCheck([]string{".md"}, md_filename) {
            c.String(http.StatusBadRequest, fmt.Sprintf("upload file is not Markdown file."))
			return
        }

        // ファイルをuploadする
		if err := c.SaveUploadedFile(md_file, PATH + md_filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload markdown file err: %s", err.Error()))
			return
        }
        cmdName := "mdtopdf"
        outputFile := strings.Replace(md_filename, filepath.Ext(md_filename), ".pdf", -1)
        if convertType == "tex" {
            cmdName = "mdtotex"
            outputFile = strings.Replace(md_filename, filepath.Ext(md_filename), ".tex", -1)
        }
        c_err := exec.Command(cmdName, PATH + md_filename).Run()
        if c_err != nil {
            log.Fatal(c_err)
        }
        fmt.Println("success!!!")
        c.Redirect(303, "http://0.0.0.0:8000/download/" + outputFile)
            
        // c.String(http.StatusOK, fmt.Sprintf("File %s waiting....", file.Filename))
    })

    router.GET("/download/:filename", func (c *gin.Context) {
    fileName := c.Param("filename")
    // fileName := "example.pdf"
    targetPath := filepath.Join(PATH, fileName)
    //This ckeck is for example, I not sure is it can prevent all possible filename attacks - will be much better if real filename will not come from user side. I not even tryed this code
    if !strings.HasPrefix(filepath.Clean(targetPath), PATH) {
        c.String(403, "Look like you attacking me")
        return
    }
    if filepath.Ext(fileName) == ".tex" {
        c.Header("Content-Type", "text/plain")
    } else {
        c.Header("Content-Type", "application/pdf")
        c.Header("Content-Disposition", "inline; filename="+fileName )
    }
    c.File(targetPath)
    })

    router.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
