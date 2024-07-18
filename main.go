package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gitsang/notation/pkg/soundslice"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type File struct {
	Href          string `json:"href,omitempty"`
	Name          string `json:"name,omitempty"`
	EnablePreview bool   `json:"enablePreview,omitempty"`
}

type Browser struct {
	Files []File `json:"files,omitempty"`
}

func (b *Browser) ToMap() map[string]any {
	jsonBytes, _ := json.Marshal(b)
	var m map[string]any
	_ = json.Unmarshal(jsonBytes, &m)
	return m
}

type Notation struct {
	Title string
	URL   string
}

func (n *Notation) ToMap() map[string]any {
	jsonBytes, _ := json.Marshal(n)
	var m map[string]any
	_ = json.Unmarshal(jsonBytes, &m)
	return m
}

func BrowserHandler(w http.ResponseWriter, r *http.Request) {
	var (
		urlpath  = r.URL.Path
		filepath = path.Join(".", urlpath)
		logger   = slog.Default().With(
			slog.String("urlpath", urlpath),
			slog.String("filepath", filepath),
		)
	)
	defer func() {
		logger.Info("Serving")
	}()

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if fileInfo.IsDir() {
		entries, err := os.ReadDir(filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		browser := Browser{
			Files: []File{},
		}
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() {
				name += "/"
			}
			browser.Files = append(browser.Files, File{
				Href:          path.Join(urlpath, entry.Name()),
				Name:          name,
				EnablePreview: path.Ext(name) == ".gp",
			})
		}
		logger = logger.With(slog.Any("browser", browser))
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, browser.ToMap())
		return
	}

	http.ServeFile(w, r, filepath)
}

func PreviewHandler(w http.ResponseWriter, r *http.Request) {
	var (
		urlpath  = r.URL.Path
		filepath = path.Join(".", strings.TrimPrefix(urlpath, "/preview"))
		logger   = slog.Default().With(
			slog.String("urlpath", urlpath),
			slog.String("filepath", filepath),
		)
	)
	defer func() {
		logger.Info("Serving")
	}()

	_ = client.DeleteNotation(r.Context(), lastSliceId)

	sliceId, err := client.CreateNotation(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger = logger.With(slog.String("sliceId", sliceId))
	lastSliceId = sliceId

	fh, err := os.Open(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fh.Close()

	_, err = client.UploadNotation(r.Context(), sliceId, fh)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scoreId, err := client.GetScoreId(sliceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = client.EnableEmbed(scoreId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notation := Notation{
		Title: strings.TrimSuffix(path.Base(filepath), ".gp"),
		URL:   fmt.Sprintf("https://www.soundslice.com/slices/%s/embed/", sliceId),
	}
	t, _ := template.ParseFiles("notation.html")
	t.Execute(w, notation.ToMap())
}

var (
	lastSliceId string
)

var client *soundslice.Client

func main() {
	godotenv.Load(".env")
	sesn := os.Getenv("SESN")
	client = soundslice.NewClient(
		soundslice.WithSesn(sesn),
		soundslice.WithLogHandler(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			})),
	)

	router := mux.NewRouter()
	router.PathPrefix("/preview").HandlerFunc(PreviewHandler)
	router.PathPrefix("/css").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
	router.PathPrefix("/js").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("js"))))
	router.NotFoundHandler = http.HandlerFunc(BrowserHandler)

	slog.Info("Listening on port :8080")
	http.ListenAndServe(":8080", router)
}
