// Package models contient les modèles et DTO du domaine d'authentification.
package models

import (
	"time"

	"comores-event-api/internal/constants"
)

// User correspond à la table users (+ rôles chargés séparément).
type User struct {
	ID               string     `db:"id"`
	Email            string     `db:"email"`
	PasswordHash     string     `db:"password_hash"`
	DisplayName      *string    `db:"display_name"`
	AccountStatus    string     `db:"account_status"`
	EmailVerifiedAt  *time.Time `db:"email_verified_at"`
	TokenVersion     int        `db:"token_version"`
	FailedLoginCount int        `db:"failed_login_count"`
	LockedUntil      *time.Time `db:"locked_until"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`

	Roles []string `db:"-"`
}

// IsActive : compte actif et non supprimé.
func (u *User) IsActive() bool {
	return u.DeletedAt == nil && u.AccountStatus == constants.StatusActive
}

// IsLocked : verrouillage temporaire en cours.
func (u *User) IsLocked(now time.Time) bool {
	return u.LockedUntil != nil && u.LockedUntil.After(now)
}

// PrimaryRole retourne le rôle principal (admin prime sur user).
func (u *User) PrimaryRole() string {
	return PrimaryRole(u.Roles)
}

// PrimaryRole choisit le rôle principal parmi une liste.
func PrimaryRole(roles []string) string {
	for _, r := range roles {
		if r == constants.RoleAdmin {
			return constants.RoleAdmin
		}
	}
	return constants.RoleUser
}

// UserResponse est l'objet utilisateur exposé par l'API.
type UserResponse struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	DisplayName   *string   `json:"display_name"`
	Role          string    `json:"role"`
	Roles         []string  `json:"roles"`
	AccountStatus string    `json:"account_status"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

// ToResponse convertit un User en réponse API (sans données sensibles).
func (u *User) ToResponse() UserResponse {
	roles := u.Roles
	if roles == nil {
		roles = []string{}
	}
	return UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		DisplayName:   u.DisplayName,
		Role:          u.PrimaryRole(),
		Roles:         roles,
		AccountStatus: u.AccountStatus,
		EmailVerified: u.EmailVerifiedAt != nil,
		CreatedAt:     u.CreatedAt,
	}
}

// RegisterRequest : corps de POST /api/auth/register.
type RegisterRequest struct {
	Email       string  `json:"email" validate:"required,email,max=254"`
	Password    string  `json:"password" validate:"required,min=8,max=72"`
	DisplayName *string `json:"display_name" validate:"omitempty,max=100"`
}

// LoginRequest : corps de POST /api/auth/login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,max=72"`
}

// UpdateProfileRequest : corps de PATCH /api/auth/me.
type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name" validate:"omitempty,max=100"`
}
