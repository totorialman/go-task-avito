// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/totorialman/go-task-avito/restapi/operations"

)

//go:generate swagger generate server --target ..\..\go-task-avito --name BackendService --spec ..\swagger.yaml --principal interface{}

func configureFlags(api *operations.BackendServiceAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.BackendServiceAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "Authorization" header is set
	if api.BearerAuthAuth == nil {
		api.BearerAuthAuth = func(token string) (interface{}, error) {
			return nil, errors.NotImplemented("api key auth (bearerAuth) Authorization from header param [Authorization] has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	if api.GetPvzHandler == nil {
		api.GetPvzHandler = operations.GetPvzHandlerFunc(func(params operations.GetPvzParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetPvz has not yet been implemented")
		})
	}
	if api.PostDummyLoginHandler == nil {
		api.PostDummyLoginHandler = operations.PostDummyLoginHandlerFunc(func(params operations.PostDummyLoginParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostDummyLogin has not yet been implemented")
		})
	}
	if api.PostLoginHandler == nil {
		api.PostLoginHandler = operations.PostLoginHandlerFunc(func(params operations.PostLoginParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostLogin has not yet been implemented")
		})
	}
	if api.PostProductsHandler == nil {
		api.PostProductsHandler = operations.PostProductsHandlerFunc(func(params operations.PostProductsParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostProducts has not yet been implemented")
		})
	}
	if api.PostPvzHandler == nil {
		api.PostPvzHandler = operations.PostPvzHandlerFunc(func(params operations.PostPvzParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostPvz has not yet been implemented")
		})
	}
	if api.PostPvzPvzIDCloseLastReceptionHandler == nil {
		api.PostPvzPvzIDCloseLastReceptionHandler = operations.PostPvzPvzIDCloseLastReceptionHandlerFunc(func(params operations.PostPvzPvzIDCloseLastReceptionParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostPvzPvzIDCloseLastReception has not yet been implemented")
		})
	}
	if api.PostPvzPvzIDDeleteLastProductHandler == nil {
		api.PostPvzPvzIDDeleteLastProductHandler = operations.PostPvzPvzIDDeleteLastProductHandlerFunc(func(params operations.PostPvzPvzIDDeleteLastProductParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostPvzPvzIDDeleteLastProduct has not yet been implemented")
		})
	}
	if api.PostReceptionsHandler == nil {
		api.PostReceptionsHandler = operations.PostReceptionsHandlerFunc(func(params operations.PostReceptionsParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostReceptions has not yet been implemented")
		})
	}
	if api.PostRegisterHandler == nil {
		api.PostRegisterHandler = operations.PostRegisterHandlerFunc(func(params operations.PostRegisterParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostRegister has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
