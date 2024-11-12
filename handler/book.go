package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/train-do/Golang-Restfull-API/model"
	"github.com/train-do/Golang-Restfull-API/service"
	"go.uber.org/zap"
)

type BookHandler struct {
	service *service.BookService
	log     *zap.Logger
}

func NewBookHandler(service *service.BookService,
	log *zap.Logger) *BookHandler {
	return &BookHandler{service, log}
}

func (h *BookHandler) GetAllBook(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Request " + r.URL.Path)
	books, err := h.service.GetAllBook()
	data := struct {
		Books []model.Book
	}{
		Books: books,
	}
	// fmt.Printf("%+v\n", data)
	if err != nil {
		log.Fatalf("Error GetAllBook: %v", err)
	}
	templates, err := template.ParseGlob("template/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
	err = templates.ExecuteTemplate(w, "book-list.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// fmt.Println("CREATE BOOK GET")
		templates, err := template.ParseGlob("template/*.html")
		if err != nil {
			log.Fatalf("Error parsing templates: %v", err)
		}
		err = templates.ExecuteTemplate(w, "add-book.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		// fmt.Println("CREATE BOOK POST")
		price, _ := strconv.Atoi(r.FormValue("price"))
		discount, _ := strconv.Atoi(r.FormValue("discount"))
		book := model.Book{
			Title:    r.FormValue("bookName"),
			Type:     r.FormValue("bookType"),
			Author:   r.FormValue("author"),
			Price:    price,
			Discount: discount,
		}

		h.service.CreateBook(&book)
		fmt.Printf("%+v+++++\n", book)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		file1, handler1, err := r.FormFile("cover")
		if err != nil {
			http.Error(w, "Error retrieving file1", http.StatusBadRequest)
			return
		}
		ext := filepath.Ext(handler1.Filename)
		handler1.Filename = book.BookCode + "cover" + ext
		defer file1.Close()

		file2, handler2, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file2", http.StatusBadRequest)
			return
		}
		ext = filepath.Ext(handler2.Filename)
		handler2.Filename = book.BookCode + "book" + ext
		defer file2.Close()

		out1, err := os.Create(filepath.Join("uploads", handler1.Filename))
		if err != nil {
			http.Error(w, "Error saving file1", http.StatusInternalServerError)
			return
		}
		defer out1.Close()

		_, err = io.Copy(out1, file1)
		if err != nil {
			http.Error(w, "Error saving file1", http.StatusInternalServerError)
			return
		}

		out2, err := os.Create(filepath.Join("uploads", handler2.Filename))
		if err != nil {
			http.Error(w, "Error saving file2", http.StatusInternalServerError)
			return
		}
		defer out2.Close()

		_, err = io.Copy(out2, file2)
		if err != nil {
			http.Error(w, "Error saving file2", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {}

func (h *BookHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("DASHBOARD")
	templates, err := template.ParseGlob("template/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
	err = templates.ExecuteTemplate(w, "dashboard.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
