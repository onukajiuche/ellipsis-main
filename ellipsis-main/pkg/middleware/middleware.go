package middleware

import (
	"brief/internal/constant"
	"brief/internal/model"
	"brief/utility"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// Admin is the middleware for admin-only endpoints
func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := getToken(r)
		if token == "" {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, "no token specified", nil)
			res, _ := json.Marshal(rd)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(res)
			return
		}

		claims, err := VerifyToken(token)
		if err != nil {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, err.Error(), nil)
			res, _ := json.Marshal(rd)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(res)
			return
		}

		// Check if request is from an admin
		if claims.Role != constant.Roles[constant.Admin] {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, "cannot access this endpoint", nil)
			res, _ := json.Marshal(rd)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(res)
			return
		}

		// Set details from token into context and execute next handler
		ctx := r.Context()
		ctx = context.WithValue(ctx, struct{}{}, &model.ContextInfo{
			ID:    claims.ID,
			Role:  claims.Role,
			Email: claims.Email,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Me is the middleware for user endpoints
func Me(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := getToken(r)
		if token == "" {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, "no token specified", nil)
			res, _ := json.Marshal(rd)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(res)
			return
		}

		claims, err := VerifyToken(token)
		if err != nil {
			rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
				constant.ErrUnauthorized, err.Error(), nil)
			res, _ := json.Marshal(rd)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(res)
			return
		}

		// Set details from token in context and execute next handler
		ctx := r.Context()
		ctx = context.WithValue(ctx, struct{}{}, &model.ContextInfo{
			ID:    claims.ID,
			Role:  claims.Role,
			Email: claims.Email,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Shorten is the middleware for the endpoint to shorten a url
// /api/v1/url/shorten - POST
func Shorten(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get token if available
		token := getToken(r)
		if token != "" {
			claims, err := VerifyToken(token)
			if err != nil {
				rd := utility.BuildErrorResponse(http.StatusUnauthorized, constant.StatusFailed,
					constant.ErrUnauthorized, err.Error(), nil)
				res, _ := json.Marshal(rd)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(res)
				return
			}

			// Set details from token in context and execute next handler
			ctx = context.WithValue(ctx, struct{}{}, &model.ContextInfo{
				ID:    claims.ID,
				Role:  claims.Role,
				Email: claims.Email,
			})

		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

// getToken contains logic to fetch token from headers
func getToken(r *http.Request) (token string) {
	auth := r.Header.Get("Authorization")
	hToken := r.Header.Get("Token") // header token
	if auth == "" {
		if hToken == "" {
			return
		} else {
			token = hToken
		}
	} else {
		// Split Authorization to get bearer token
		strs := strings.Split(auth, " ")
		if len(strs) > 1 {
			token = strs[1]
		} else {
			token = auth
		}
	}

	return
}
