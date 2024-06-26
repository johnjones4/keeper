// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
)

// Message defines model for Message.
type Message struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

// Note defines model for Note.
type Note struct {
	Body     string     `json:"body"`
	Key      string     `json:"key"`
	Modified *time.Time `json:"modified,omitempty"`
}

// Notes defines model for Notes.
type Notes struct {
	NextPage *string  `json:"nextPage,omitempty"`
	Notes    []string `json:"notes"`
}

// TokenRequest defines model for TokenRequest.
type TokenRequest struct {
	Password string `json:"password"`
}

// TokenResponse defines model for TokenResponse.
type TokenResponse struct {
	Token string `json:"token"`
}

// GetNoteParams defines parameters for GetNote.
type GetNoteParams struct {
	Q             *string `form:"q,omitempty" json:"q,omitempty"`
	Dir           *string `form:"dir,omitempty" json:"dir,omitempty"`
	Page          *string `form:"page,omitempty" json:"page,omitempty"`
	Authorization string  `json:"Authorization"`
}

// PostNoteParams defines parameters for PostNote.
type PostNoteParams struct {
	Authorization string `json:"Authorization"`
}

// GetNoteKeyParams defines parameters for GetNoteKey.
type GetNoteKeyParams struct {
	Authorization string `json:"Authorization"`
}

// PutNoteKeyParams defines parameters for PutNoteKey.
type PutNoteKeyParams struct {
	Authorization string `json:"Authorization"`
}

// PostNoteJSONRequestBody defines body for PostNote for application/json ContentType.
type PostNoteJSONRequestBody = Note

// PutNoteKeyJSONRequestBody defines body for PutNoteKey for application/json ContentType.
type PutNoteKeyJSONRequestBody = Note

// PostTokenJSONRequestBody defines body for PostToken for application/json ContentType.
type PostTokenJSONRequestBody = TokenRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /health)
	GetHealth(w http.ResponseWriter, r *http.Request)

	// (GET /note)
	GetNote(w http.ResponseWriter, r *http.Request, params GetNoteParams)

	// (POST /note)
	PostNote(w http.ResponseWriter, r *http.Request, params PostNoteParams)

	// (GET /note/{key})
	GetNoteKey(w http.ResponseWriter, r *http.Request, key string, params GetNoteKeyParams)

	// (PUT /note/{key})
	PutNoteKey(w http.ResponseWriter, r *http.Request, key string, params PutNoteKeyParams)

	// (POST /token)
	PostToken(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (GET /health)
func (_ Unimplemented) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /note)
func (_ Unimplemented) GetNote(w http.ResponseWriter, r *http.Request, params GetNoteParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /note)
func (_ Unimplemented) PostNote(w http.ResponseWriter, r *http.Request, params PostNoteParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /note/{key})
func (_ Unimplemented) GetNoteKey(w http.ResponseWriter, r *http.Request, key string, params GetNoteKeyParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (PUT /note/{key})
func (_ Unimplemented) PutNoteKey(w http.ResponseWriter, r *http.Request, key string, params PutNoteKeyParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /token)
func (_ Unimplemented) PostToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetHealth operation middleware
func (siw *ServerInterfaceWrapper) GetHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHealth(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetNote operation middleware
func (siw *ServerInterfaceWrapper) GetNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetNoteParams

	// ------------- Optional query parameter "q" -------------

	err = runtime.BindQueryParameter("form", true, false, "q", r.URL.Query(), &params.Q)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "q", Err: err})
		return
	}

	// ------------- Optional query parameter "dir" -------------

	err = runtime.BindQueryParameter("form", true, false, "dir", r.URL.Query(), &params.Dir)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "dir", Err: err})
		return
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page", Err: err})
		return
	}

	headers := r.Header

	// ------------- Required header parameter "Authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Authorization")]; found {
		var Authorization string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Authorization", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "Authorization", runtime.ParamLocationHeader, valueList[0], &Authorization)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Authorization", Err: err})
			return
		}

		params.Authorization = Authorization

	} else {
		err := fmt.Errorf("Header parameter Authorization is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "Authorization", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetNote(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostNote operation middleware
func (siw *ServerInterfaceWrapper) PostNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostNoteParams

	headers := r.Header

	// ------------- Required header parameter "Authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Authorization")]; found {
		var Authorization string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Authorization", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "Authorization", runtime.ParamLocationHeader, valueList[0], &Authorization)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Authorization", Err: err})
			return
		}

		params.Authorization = Authorization

	} else {
		err := fmt.Errorf("Header parameter Authorization is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "Authorization", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostNote(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetNoteKey operation middleware
func (siw *ServerInterfaceWrapper) GetNoteKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "key" -------------
	var key string

	err = runtime.BindStyledParameterWithLocation("simple", false, "key", runtime.ParamLocationPath, chi.URLParam(r, "key"), &key)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "key", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetNoteKeyParams

	headers := r.Header

	// ------------- Required header parameter "Authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Authorization")]; found {
		var Authorization string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Authorization", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "Authorization", runtime.ParamLocationHeader, valueList[0], &Authorization)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Authorization", Err: err})
			return
		}

		params.Authorization = Authorization

	} else {
		err := fmt.Errorf("Header parameter Authorization is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "Authorization", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetNoteKey(w, r, key, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PutNoteKey operation middleware
func (siw *ServerInterfaceWrapper) PutNoteKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "key" -------------
	var key string

	err = runtime.BindStyledParameterWithLocation("simple", false, "key", runtime.ParamLocationPath, chi.URLParam(r, "key"), &key)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "key", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PutNoteKeyParams

	headers := r.Header

	// ------------- Required header parameter "Authorization" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Authorization")]; found {
		var Authorization string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "Authorization", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "Authorization", runtime.ParamLocationHeader, valueList[0], &Authorization)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Authorization", Err: err})
			return
		}

		params.Authorization = Authorization

	} else {
		err := fmt.Errorf("Header parameter Authorization is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "Authorization", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PutNoteKey(w, r, key, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostToken operation middleware
func (siw *ServerInterfaceWrapper) PostToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostToken(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/health", wrapper.GetHealth)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/note", wrapper.GetNote)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/note", wrapper.PostNote)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/note/{key}", wrapper.GetNoteKey)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/note/{key}", wrapper.PutNoteKey)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/token", wrapper.PostToken)
	})

	return r
}
