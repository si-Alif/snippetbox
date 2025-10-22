package main

// defining a contextKey type so that it doesn't make any collision with other keys
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
