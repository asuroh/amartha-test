package handler

import (
	"amartha-test/pkg/str"
	"amartha-test/usecase"
	"encoding/json"
	"net/http"
	"strings"

	ut "github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
)

// Handler ...
type Handler struct {
	ContractUC *usecase.ContractUC
	EnvConfig  map[string]string
	Validate   *validator.Validate
	Translator ut.Translator
}

type Error struct {
	Error interface{} `json:"error"`
}

// Bind bind the API request payload (body) into request struct.
func (h Handler) Bind(r *http.Request, input interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&input)

	return err
}

// SendSuccess send success into response with 200 http code.
func SendSuccess(w http.ResponseWriter, payload interface{}) {
	RespondWithJSON(w, 200, "Success", payload)
}

// SendBadRequest send bad request into response with 400 http code.
func SendBadRequest(w http.ResponseWriter, payload string) {
	RespondWithJSON(w, 400, "fail", Error{Error: payload})
}

// SendNotFound send bad request into response with 404 http code.
func SendNotFound(w http.ResponseWriter, payload string) {
	RespondWithJSON(w, 404, "fail", Error{Error: "Data Not Found"})
}

// SendRequestValidationError Send validation error response to consumers.
func (h Handler) SendRequestValidationError(w http.ResponseWriter, validationErrors validator.ValidationErrors) {
	errorResponse := map[string][]string{}
	errorTranslation := validationErrors.Translate(h.Translator)
	for _, err := range validationErrors {
		errKey := str.Underscore(err.StructField())
		errorResponse[errKey] = append(
			errorResponse[errKey],
			strings.Replace(errorTranslation[err.Namespace()], err.StructField(), "[]", -1),
		)
	}

	RespondWithJSON(w, 400, "fail", Error{Error: errorResponse})
}

// RespondWithJSON write json response format
func RespondWithJSON(w http.ResponseWriter, httpCode int, message string, payload interface{}) {
	respPayload := map[string]interface{}{
		"status": message,
		"data":   payload,
	}

	response, _ := json.Marshal(respPayload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(response)
}
