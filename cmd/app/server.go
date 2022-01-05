package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Tursunkhuja/http/pkg/banners"
)

//it is logical server of our app
type Server struct {
	mux        *http.ServeMux
	bannersSvc *banners.Service
}

func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: bannersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

//initializing server
func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleremoveByID)
}

// get All
func (s *Server) handleGetAllBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := s.bannersSvc.All(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}

	data, err := json.Marshal(banners)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

//get by ID
func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

// Save and update
/*
func (s *Server) handleSaveBanner(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	titleParam := r.URL.Query().Get("title")
	contentParam := r.URL.Query().Get("content")
	buttonParam := r.URL.Query().Get("button")
	linkParam := r.URL.Query().Get("link")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	banner := banners.Banner{
		ID:      id,
		Title:   titleParam,
		Content: contentParam,
		Button:  buttonParam,
		Link:    linkParam,
	}

	GotB, err := s.bannersSvc.Save(r.Context(), &banner)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(GotB)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Contetn-Type", "applicatrion/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}

}
*/

// delete banner byID
func (s *Server) handleremoveByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	dBanner, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(dBanner)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Contetn-Type", "applicatrion/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleSaveBanner(w http.ResponseWriter, r *http.Request) {
	idp := r.PostFormValue("id")
	title := r.PostFormValue("title")
	content := r.PostFormValue("content")
	button := r.PostFormValue("button")
	link := r.PostFormValue("link")

	id, err := strconv.ParseInt(idp, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if title == "" && content == "" && button == "" && link == "" {

		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item := &banners.Banner{
		ID:      id,
		Title:   title,
		Content: content,
		Button:  button,
		Link:    link,
	}

	file, header, err := r.FormFile("image")

	if err == nil {
		name := strings.Split(header.Filename, ".")
		item.Image = name[len(name)-1]
	}

	banner, err := s.bannersSvc.Save(r.Context(), item, file)

	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(banner)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}
