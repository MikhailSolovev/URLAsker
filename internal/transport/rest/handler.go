// Package classification URLASker.
//
//	Schemes: http
//	BasePath: /api/v1
//	Version: 0.0.1
//	Host: localhost
//
// swagger:meta
package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MikhailSolovev/URLAsker/internal/interfaces"
	"github.com/MikhailSolovev/URLAsker/internal/models"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"time"
)

// TODO: prometheus metrics (another contour - port)
// TODO: give more explanation in user error response

type Handler struct {
	router *mux.Router
	asker  interfaces.Asker
}

func New(router *mux.Router, asker interfaces.Asker) *Handler {
	return &Handler{router: router, asker: asker}
}

func (h *Handler) Register() {
	h.router.HandleFunc("/readiness", h.Alive).Methods(http.MethodGet)
	h.router.HandleFunc("/liveness", h.Alive).Methods(http.MethodGet)
	apiRouter := h.router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(SetCORSHeaders)
	apiRouter.HandleFunc("/info", HandleError(h.GetInfo)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/listLatest", HandleError(h.ListLatestResult)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/list", HandleError(h.ListResults)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/setInterval", HandleError(h.SetInterval)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/setUrls", HandleError(h.SetURLs)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/addUrls", HandleError(h.AddURLs)).Methods(http.MethodPut)
	apiRouter.HandleFunc("/deleteUrls", HandleError(h.DeleteURLs)).Methods(http.MethodDelete)
	apiRouter.Use(mux.CORSMethodMiddleware(h.router))
}

// swagger:operation GET /info Methods idInfo
//
// Get info about asker service.
// ---
// responses:
//
//		'200': infoResponse
//		'500':
//	    description: internal server error
//	    type: string
func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) error {
	data, err := h.asker.GetInfo(context.Background())
	if err != nil {
		return InternalServerErr.SetDebugMsg(fmt.Sprintf("failed to get info due to error: %v",
			err.Error()))
	}

	var resp models.InfoRestDTO
	resp.Interval = data.Interval.String()
	for key := range data.URLs {
		resp.URLs = append(resp.URLs, key)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(resp)
	if err != nil {
		return InternalServerErr.SetDebugMsg(fmt.Sprintf("failed to marshall info to response due to error: %v",
			err.Error()))
	}
	w.Write(body)

	return nil
}

// swagger:operation GET /listLatest Methods idListLatest
//
// Get latest result.
// ---
// responses:
//
//		'200': listLatestResponse
//	 	'404':
//	    description: not found
//	    type: string
//		'500':
//	    description: internal server error
//	    type: string
func (h *Handler) ListLatestResult(w http.ResponseWriter, r *http.Request) error {
	data, err := h.asker.ListLatestResult(context.Background())
	if err != nil {
		return InternalServerErr.SetDebugMsg(fmt.Sprintf("failed to get result due to error: %v",
			err.Error()))
	}

	if len(data.URLs) == 0 {
		return NotFoundErr.SetDebugMsg("not found")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(data)
	if err != nil {
		return InternalServerErr.SetDebugMsg(fmt.Sprintf("failed to marshall result to response due to error: %v",
			err.Error()))
	}
	w.Write(body)

	return nil
}

// swagger:operation GET /list Methods idList
//
// Get all results between two dates.
// ---
// responses:
//
//		'200': listResponse
//	 	'404':
//	    description: not found
//	    type: string
//		'500':
//	    description: internal server error
//	    type: string
func (h *Handler) ListResults(w http.ResponseWriter, r *http.Request) error {
	// RFC3339 time format, RFC3339: year-month-day T hours-minutes-seconds Z
	dateFromStr := r.URL.Query().Get("dateFrom")
	dateFrom, err := time.Parse(time.RFC3339, dateFromStr)
	if err != nil {
		return BadReqErr.SetDebugMsg(fmt.Sprintf("failed to parse dateFrom param due to error: %v", err.Error()))
	}
	// RFC3339 time format, RFC3339: year-month-day T hours-minutes-seconds Z
	dateToStr := r.URL.Query().Get("dateTo")
	dateTo, err := time.Parse(time.RFC3339, dateToStr)
	if err != nil {
		return BadReqErr.SetDebugMsg(fmt.Sprintf("failed to parse dateTo param due to error: %v", err.Error()))
	}

	if !dateTo.After(dateFrom) {
		return BadReqErr.SetDebugMsg("dateFrom latter than dateTo param")
	}

	data, err := h.asker.ListResults(context.Background(), dateFrom, dateTo)
	if err != nil {
		return InternalServerErr.SetDebugMsg(fmt.Sprintf("failed to get results due to error: %v",
			err.Error()))
	}

	if len(data.Results) == 0 {
		return NotFoundErr.SetDebugMsg("not found")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(data)
	if err != nil {
		return InternalServerErr.SetDebugMsg(fmt.Sprintf("failed to marshall results to response due to error: %v",
			err.Error()))
	}
	w.Write(body)

	return nil
}

// swagger:operation POST /setInterval Methods idSetInterval
//
// Set interval of asking.
// ---
// produces:
// - text/plain
// responses:
//
//		'200':
//	    description: Success
//	    type: string
//	 	'400':
//	    description: bad request
//	    type: string
//		'500':
//	    description: internal server error
//	    type: string
func (h *Handler) SetInterval(w http.ResponseWriter, r *http.Request) error {
	intervalStr := r.URL.Query().Get("interval")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return BadReqErr.SetDebugMsg(fmt.Sprintf("failed to parse interval param due to error: %v", err.Error()))
	}

	if err = h.asker.SetInterval(context.Background(), interval); err != nil {
		return InternalServerErr.SetDebugMsg(fmt.Sprintf("failed to set interval due to error: %v",
			err.Error()))
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))

	return nil
}

// swagger:operation POST /setUrls Methods idSetUrls
//
// Rewrite set of urls.
// ---
// produces:
// - text/plain
// responses:
//
//		'200':
//	    description: Success
//	    type: string
//	 	'400':
//	    description: bad request
//	    type: string
//		'500':
//	    description: internal server error
//	    type: string
func (h *Handler) SetURLs(w http.ResponseWriter, r *http.Request) error {
	var data []string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return BadReqErr.SetDebugMsg(fmt.Sprintf("failed to parse body due to error: %v", err.Error()))
	}
	for _, rawUrl := range data {
		if _, err = url.ParseRequestURI(rawUrl); err != nil {
			return BadReqErr.SetDebugMsg(fmt.Sprintf("invalid url in body due to error: %v", err.Error()))
		}
	}

	if err = h.asker.SetURLs(context.Background(), data...); err != nil {
		return InternalServerErr.SetDebugMsg(fmt.Sprintf("failed to set urls due to error: %v",
			err.Error()))
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Success"))

	return nil
}

// swagger:operation PUT /addUrls Methods idAddUrls
//
// Append urls to set of urls.
// ---
// produces:
// - text/plain
// responses:
//
//		'200':
//	    description: Success
//	    type: string
//	 	'400':
//	    description: bad request
//	    type: string
//		'500':
//	    description: internal server error
//	    type: string
func (h *Handler) AddURLs(w http.ResponseWriter, r *http.Request) error {
	var data []string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return BadReqErr.SetDebugMsg(fmt.Sprintf("failed to parse body due to error: %v", err.Error()))
	}
	for _, rawUrl := range data {
		if _, err = url.ParseRequestURI(rawUrl); err != nil {
			return BadReqErr.SetDebugMsg(fmt.Sprintf("invalid url in body due to error: %v", err.Error()))
		}
	}

	if err = h.asker.AddURLs(context.Background(), data...); err != nil {
		return InternalServerErr.SetDebugMsg(fmt.Sprintf("failed to add urls due to error: %v",
			err.Error()))
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))

	return nil
}

// swagger:operation DELETE /deleteUrls Methods idDeleteUrls
//
// Delete urls from set of urls.
// ---
// produces:
// - text/plain
// responses:
//
//		'200':
//	    description: Success
//	    type: string
//	 	'400':
//	    description: bad request
//	    type: string
//		'500':
//	    description: internal server error
//	    type: string
func (h *Handler) DeleteURLs(w http.ResponseWriter, r *http.Request) error {
	var data []string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return BadReqErr.SetDebugMsg(fmt.Sprintf("failed to parse body due to error: %v", err.Error()))
	}
	for _, rawUrl := range data {
		if _, err = url.ParseRequestURI(rawUrl); err != nil {
			return BadReqErr.SetDebugMsg(fmt.Sprintf("invalid url in body due to error: %v", err.Error()))
		}
	}

	if err = h.asker.DeleteURLs(context.Background(), data...); err != nil {
		return InternalServerErr.SetDebugMsg(fmt.Sprintf("failed to delete urls due to error: %v",
			err.Error()))
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))

	return nil
}

func (h *Handler) Alive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
