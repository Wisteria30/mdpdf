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

type Result struct {
    Url string `json: "url"`
}

func main() {
    // host := "localhost"
    // port := "8000" 
    router := gin.Default()
    const PATH = "resource/"
    
    router.LoadHTMLGlob("public/*.html")
    
    router.GET("/", func(c *gin.Context) {
        c.HTML(200, "index.html", gin.H{})
    })

	router.POST("/upload", func(c *gin.Context) {
        convertType := c.PostForm("type")
        // Source
		file, err := c.FormFile("file")
		if err != nil {
            c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
        
        filename := filepath.Base(file.Filename)

        // md以外をはじく
        if filepath.Ext(filename) != ".md" {
            c.String(http.StatusBadRequest, fmt.Sprintf("upload file is not Markdown file."))
			return
        }

        // ファイルをuploadする
		if err := c.SaveUploadedFile(file, PATH + filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
        }

        // go func(filename string, c *gin.Context) {
        //     cmd := exec.Command("mdtopdf", PATH + filename)
        //     stdoutStderr, err := cmd.CombinedOutput()
        //     if err != nil {
        //         log.Fatal(err)
        //     }
        //     fmt.Printf("%s", stdoutStderr)
        //     fmt.Println("success!!!")
        //     c.Redirect(303, "http://www.google.com/")
        //     }(filename, c)
        cmdName := "mdtopdf"
        outputFile := strings.Replace(filename, filepath.Ext(filename), ".pdf", -1)
        if convertType == "tex" {
            cmdName = "mdtotex"
            outputFile = strings.Replace(filename, filepath.Ext(filename), ".tex", -1)
        }
        cmd := exec.Command(cmdName, PATH + filename)
        stdoutStderr, err := cmd.CombinedOutput()
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s", stdoutStderr)
        fmt.Println("success!!!")

        result := Result{
            // Url: "http://"+host+":"+port+"/download/" + outputFile,
            Url: "/api/download/" + outputFile,
        }

        c.JSON(200, result)

        // c.Redirect(303, "http://0.0.0.0:8000/download/" + outputFile)
            
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
