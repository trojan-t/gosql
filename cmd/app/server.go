package app

import (
	"encoding/json"
	"github.com/trojan-t/gosql/pkg/customers"
	"log"
	"errors"
	"strconv"
	"net/http"
)

// Server is struct
type Server struct {
	mux         *http.ServeMux
	customerSvc *customers.Service
}

// NewServer is function constructor
func NewServer(mux *http.ServeMux, customerSvc *customers.Service) *Server {
	return &Server{mux: mux, customerSvc: customerSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Init is init method
func (s *Server) Init() {
	log.Println("Init method")
	s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
	s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomers)
	s.mux.HandleFunc("/customers.getAllActive", s.handleGetAllActiveCustomers)
	s.mux.HandleFunc("/customers.blockById", s.handleBlockByID)
	s.mux.HandleFunc("/customers.unblockById", s.handleUnBlockByID)
	s.mux.HandleFunc("/customers.removeById", s.handleDelete)
	s.mux.HandleFunc("/customers.save", s.handleSave)
}

// handleGetCustomerByID is method
func (s *Server) handleGetCustomerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.ByID(request.Context(), id)
	log.Println(item)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(writer, http.StatusNotFound, err)
		return
	}

	if err != nil {
		log.Println(err)
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(writer, item)
}

func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {
	items, err := s.customerSvc.All(r.Context())

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, items)
}

func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request) {
	items, err := s.customerSvc.AllActive(r.Context())

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, items)
}

func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.ChangeActive(r.Context(), id, false)

	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, item)
}

func (s *Server) handleUnBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.ChangeActive(r.Context(), id, true)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	jsonResponse(w, item)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	idP := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.Delete(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, item)
}

func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	idP := r.FormValue("id")
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	id, err := strconv.ParseInt(idP, 10, 64)

	if err != nil {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	if name == "" && phone == "" {
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item := &customers.Customer{ID: id, Name: name, Phone: phone}
	customer, err := s.customerSvc.Save(r.Context(), item)

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, customer)
}

func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	log.Print(err)
	http.Error(w, http.StatusText(httpSts), httpSts)
}

func jsonResponse(writer http.ResponseWriter, data interface{}) {
	item, err := json.Marshal(data)
	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(item)
	if err != nil {
		log.Println("Error write response: ", err)
	}
}