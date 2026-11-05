// Package handlers : handlers HTTP de l'API.
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func init() {
	// Utilise le nom JSON des champs dans les erreurs de validation.
	validate.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return f.Name
		}
		return name
	})
}

// WriteJSON écrit une réponse JSON.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// WriteError écrit {"error":"..."}.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

// decodeJSON lit un corps JSON (1 Mo max) dans dst et le valide.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		var maxErr *http.MaxBytesError
		switch {
		case errors.As(err, &maxErr):
			WriteError(w, http.StatusRequestEntityTooLarge, "Corps de requête trop volumineux")
		case errors.Is(err, io.EOF):
			WriteError(w, http.StatusBadRequest, "Corps de requête vide")
		default:
			WriteError(w, http.StatusBadRequest, "JSON invalide")
		}
		return false
	}
	if err := validate.Struct(dst); err != nil {
		WriteError(w, http.StatusBadRequest, validationMessage(err))
		return false
	}
	return true
}

func validationMessage(err error) string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) || len(verrs) == 0 {
		return "Données invalides"
	}
	fe := verrs[0]
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("Le champ %s est obligatoire", field)
	case "email":
		return "Adresse email invalide"
	case "min":
		return fmt.Sprintf("Le champ %s doit contenir au moins %s caractères", field, fe.Param())
	case "max":
		return fmt.Sprintf("Le champ %s doit contenir au plus %s caractères", field, fe.Param())
	default:
		return fmt.Sprintf("Le champ %s est invalide", field)
	}
}
