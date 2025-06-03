package main

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// applyCORSHandler applies a CORS policy to the router. CORS stands for Cross-Origin Resource Sharing: it's a security
// feature present in web browsers that blocks JavaScript requests going across different domains if not specified in a
// policy. This function sends the policy of this API server.
func applyCORSHandler(h http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
			"Accept-Language",
			"Accept-Encoding",
		}),
		handlers.AllowedMethods([]string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
			"HEAD",
			"PATCH",
		}),

		handlers.AllowedOrigins([]string{
			"*",
		}),

		handlers.AllowCredentials(),

		handlers.MaxAge(43200),

		handlers.ExposedHeaders([]string{
			"X-Total-Count",
			"X-Pagination",
		}),
	)(h)
}
