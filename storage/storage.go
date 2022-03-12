package storage

import (
	"musica/auth"
	"musica/db"
	"musica/model"
	"musica/utils"
	"os/exec"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"io"
	"log"
	"mime"
	"net/http"
	"os"

	"strings"
	"time"
)

var d = db.Db()
var path = utils.Env("TPATH")

type freq struct {
	Object string
	File   string
	Bucket string
}

func GetFiles(c *gin.Context) {
	claims, is_auth, message := auth.CheckToken(c)
	if !is_auth {
		c.JSON(http.StatusOK, gin.H{"message": message, "error": "INVALID_TOKEN"})
		return
	}
	var files []model.File
	d.Where("user_id = ?", fmt.Sprintf("%v", claims["user_id"])).Find(&files)
	c.JSON(http.StatusOK, files)
}
func GetFile(c *gin.Context) {
	claims, is_auth, message := auth.CheckToken(c)
	if !is_auth {
		c.JSON(http.StatusOK, gin.H{"message": message, "error": "INVALID_TOKEN"})
		return
	}
	var files []model.File
	d.Where("user_id = ?", fmt.Sprintf("%v", claims["user_id"])).Find(&files)

	c.JSON(http.StatusOK, files)
}

func Upload(c *gin.Context) {
	claims, is_auth, message := auth.CheckToken(c)
	if !is_auth {
		c.JSON(http.StatusOK, gin.H{"message": message, "error": "INVALID_TOKEN"})
		return
	}

	var ufile model.File
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	fuid := uuid.NewString()
	ufile.ID = fuid
	filename := header.Filename
	ext_sl := strings.Split(filename, ".")
	ext := "." + ext_sl[len(ext_sl)-1]
	slugged_name := ufile.ID + ext //strings.Join(strings.Split(filename, " "), "-")
	dir := utils.Path("public/") + slugged_name
	out, err := os.Create(dir)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Println(err)
	}
	fi, err := os.Stat(dir)
	if err != nil {
		log.Println(err)
	}

	ufile.OriginalName = filename
	ufile.CreatedAt = time.Now()

	ufile.MimeType = GetMimeType(ext)
	//filepath := ""//"http://localhost:8080/file/" + slugged_name

	ufile.Size = fi.Size()
	ufile.UserId = fmt.Sprintf("%v", claims["user_id"])
	//orgfiledir := path + "/public/" + slugged_name

	d.Create(&ufile)
	c.JSON(http.StatusOK, ufile)

}

func GetMimeType(fileExtension string) string {
	return mime.TypeByExtension(fileExtension)
}

func GetFingerprint(file string) string {
	out, _ := exec.Command(utils.Path("fpcalc"), file).Output()
	split := strings.Split(string(out), "FINGERPRINT=")
	return split[1]
}
