package zapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Req struct {
	*http.Request
	Params map[string]string
}

type Res struct {
	writer http.ResponseWriter
}

type HandlerFunc func(*Req, Res)

type CookieOptions struct {
	Path     string
	Domain   string
	Expires  time.Time
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

func (res *Res) JSON(v any) {
	res.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(res.writer, err.Error(), http.StatusInternalServerError)
		return
	}
	res.writer.Write(data)
}

func (req *Req) BindJSON(dest any) error {
	defer req.Body.Close()
	if !strings.HasPrefix(req.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("invalid content-type, expected application/json")
	}
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(dest)
}

func (res *Res) Header() http.Header {
	return res.writer.Header()
}

func (res *Res) Write(str string) (int, error) {
	return res.writer.Write([]byte(str))
}

func (res *Res) WriteHeader(statusCode int) {
	res.writer.WriteHeader(statusCode)
}

func (res *Res) Redirect(url string) {
	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusFound)
}

var mimeTypes = map[string]string{
	".aac":   "audio/aac",
	".avi":   "video/x-msvideo",
	".bmp":   "image/bmp",
	".css":   "text/css",
	".csv":   "text/csv",
	".doc":   "application/msword",
	".docx":  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".eot":   "application/vnd.ms-fontobject",
	".gif":   "image/gif",
	".html":  "text/html",
	".ico":   "image/vnd.microsoft.icon",
	".jpeg":  "image/jpeg",
	".jpg":   "image/jpeg",
	".js":    "application/javascript",
	".json":  "application/json",
	".mp3":   "audio/mpeg",
	".mpeg":  "video/mpeg",
	".mp4":   "video/mp4",
	".otf":   "font/otf",
	".pdf":   "application/pdf",
	".png":   "image/png",
	".svg":   "image/svg+xml",
	".ttf":   "font/ttf",
	".txt":   "text/plain",
	".wav":   "audio/wav",
	".woff":  "font/woff",
	".woff2": "font/woff2",
	".xls":   "application/vnd.ms-excel",
	".xlsx":  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".xml":   "application/xml",
	".zip":   "application/zip",
}

func (res *Res) SendFile(filePath string) error {

	data, err := os.ReadFile(filePath)

	if err != nil {
		http.Error(res.writer, "File not found", http.StatusNotFound)
		return err
	}

	ext := filepath.Ext(filePath)
	mime, exists := mimeTypes[ext]
	if !exists {
		mime = "application/octet-stream"
	}

	if ext != "" && mime != "" {
		res.Header().Set("Content-Type", mime)
	}

	res.writer.Write(data)
	return nil
}

func (req *Req) GetCookie(name string) (string, error) {
	cookie, err := req.Request.Cookie(name)
	if err != nil {
		return "", err // returns http.ErrNoCookie if not found
	}
	return cookie.Value, nil
}

func (res *Res) SetCookie(name, value string, opts *CookieOptions) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}

	if opts != nil {
		if opts.Path != "" {
			cookie.Path = opts.Path
		}
		if opts.Domain != "" {
			cookie.Domain = opts.Domain
		}
		if !opts.Expires.IsZero() {
			cookie.Expires = opts.Expires
		}
		if opts.MaxAge != 0 {
			cookie.MaxAge = opts.MaxAge
		}
		cookie.Secure = opts.Secure
		cookie.HttpOnly = opts.HttpOnly
		cookie.SameSite = opts.SameSite
	}

	http.SetCookie(res.writer, cookie)
}

func (res *Res) DeleteCookie(name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(res.writer, cookie)
}
