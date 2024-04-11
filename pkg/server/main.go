package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Charlie-Root/npv/pkg/db"
	"github.com/Charlie-Root/npv/pkg/logging"
)

var logger = logging.NewLogger("parser")

// Handler represents the API handler with a reference to the database.
type Handler struct {
	DB *db.DB
}

// NewHandler initializes a new API handler.
func NewHandler(database *db.DB) *Handler {
	return &Handler{DB: database}
}

// AddHostHandler handles requests to add a host to the database.
func (h *Handler) AddHostHandler(w http.ResponseWriter, r *http.Request) {
	var host db.Host
	err := json.NewDecoder(r.Body).Decode(&host)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	logger.Warning("Checking if hosts exists")

	// Check if the host already exists
	existingHost, _, existingHostCount := h.DB.GetHostByAddress(host.Address, host.HostTTL)
	if existingHost != "" {
		// If the host already exists, update its count
		host.HostCount += existingHostCount
		h.DB.UpdateHostCount(host.Address, existingHostCount)
		logger.Warning("Host exists, updating count")
		w.WriteHeader(http.StatusCreated)
		return
	}
	logger.Warning("Host does not exist, inserting")

	err = h.DB.InsertHost(host)
	if err != nil {
		http.Error(w, "Failed to insert host", http.StatusInternalServerError)
		return
	}
	logger.Warning("Host is inserted in the DB")

	w.WriteHeader(http.StatusCreated)
}
// AddLinkHandler handles requests to add a link to the database.
func (h *Handler) AddLinkHandler(w http.ResponseWriter, r *http.Request) {
	var link db.Link
	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the link already exists
	existingLink, _, _, linkCount := h.DB.GetLink(link.Source, link.Target, link.TargetTTL)
	if existingLink != "" {
		// If the link already exists, update its count
		link.LinkCount += linkCount
		h.DB.UpdateLinkCount(link.Source, link.Target, link.TargetTTL, linkCount)
		w.WriteHeader(http.StatusCreated)
		return
	}

	err = h.DB.InsertLink(link)
	if err != nil {
		http.Error(w, "Failed to insert link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) HostExistsHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request parameters
    query := r.URL.Query()
    hostname := query.Get("hostname")
    ttlStr := query.Get("ttl")
    ttl, err := strconv.Atoi(ttlStr)
    if err != nil {
        http.Error(w, "Invalid TTL", http.StatusBadRequest)
        return
    }
	logger.Warning(fmt.Sprintf("Looked up host %s with TTL %d", hostname, ttl))

    // Check if the host exists in the database
    exists, _, _ := h.DB.GetHostByAddress(hostname, ttl)

	logger.Warning(fmt.Sprintf("Host %s with TTL %d exists: %s", hostname, ttl, exists))
    // Send response
    json.NewEncoder(w).Encode(map[string]bool{"exists": exists == "true"})
}
func (h *Handler) LinkExistsHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request parameters
    query := r.URL.Query()
    source := query.Get("source")
    target := query.Get("target")
    ttlStr := query.Get("ttl")
    ttl, err := strconv.Atoi(ttlStr)
    if err != nil {
        http.Error(w, "Invalid TTL", http.StatusBadRequest)
        return
    }

    // Check if the link exists in the database
    exists, _, _, _ := h.DB.GetLink(source, target, ttl)

    // Send response
    json.NewEncoder(w).Encode(map[string]bool{"exists": exists == "true"})
}


// SetupRoutes sets up API routes.
func (h *Handler) SetupRoutes() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/api/add-host", h.AddHostHandler)
	r.HandleFunc("/api/add-link", h.AddLinkHandler)
	r.HandleFunc("/api/host-exists", h.HostExistsHandler)
	r.HandleFunc("/api/link-exists", h.LinkExistsHandler)
	

	return r
}
