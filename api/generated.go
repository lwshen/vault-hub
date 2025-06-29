// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// GetUserResponse defines model for GetUserResponse.
type GetUserResponse struct {
	Avatar *string             `json:"avatar,omitempty"`
	Email  openapi_types.Email `json:"email"`
	Name   *string             `json:"name,omitempty"`
}

// HealthCheckResponse defines model for HealthCheckResponse.
type HealthCheckResponse struct {
	Status    *string    `json:"status,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Email    openapi_types.Email `json:"email"`
	Password string              `json:"password"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Token string `json:"token"`
}

// SignupRequest defines model for SignupRequest.
type SignupRequest struct {
	Email    openapi_types.Email `json:"email"`
	Name     string              `json:"name"`
	Password string              `json:"password"`
}

// SignupResponse defines model for SignupResponse.
type SignupResponse struct {
	Token string `json:"token"`
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = LoginRequest

// SignupJSONRequestBody defines body for Signup for application/json ContentType.
type SignupJSONRequestBody = SignupRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /api/auth/login)
	Login(c *fiber.Ctx) error

	// (GET /api/auth/logout)
	Logout(c *fiber.Ctx) error

	// (POST /api/auth/signup)
	Signup(c *fiber.Ctx) error

	// (GET /api/health)
	Health(c *fiber.Ctx) error

	// (GET /api/user)
	GetCurrentUser(c *fiber.Ctx) error
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

type MiddlewareFunc fiber.Handler

// Login operation middleware
func (siw *ServerInterfaceWrapper) Login(c *fiber.Ctx) error {

	return siw.Handler.Login(c)
}

// Logout operation middleware
func (siw *ServerInterfaceWrapper) Logout(c *fiber.Ctx) error {

	return siw.Handler.Logout(c)
}

// Signup operation middleware
func (siw *ServerInterfaceWrapper) Signup(c *fiber.Ctx) error {

	return siw.Handler.Signup(c)
}

// Health operation middleware
func (siw *ServerInterfaceWrapper) Health(c *fiber.Ctx) error {

	return siw.Handler.Health(c)
}

// GetCurrentUser operation middleware
func (siw *ServerInterfaceWrapper) GetCurrentUser(c *fiber.Ctx) error {

	return siw.Handler.GetCurrentUser(c)
}

// FiberServerOptions provides options for the Fiber server.
type FiberServerOptions struct {
	BaseURL     string
	Middlewares []MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router fiber.Router, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, FiberServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router fiber.Router, si ServerInterface, options FiberServerOptions) {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	for _, m := range options.Middlewares {
		router.Use(fiber.Handler(m))
	}

	router.Post(options.BaseURL+"/api/auth/login", wrapper.Login)

	router.Get(options.BaseURL+"/api/auth/logout", wrapper.Logout)

	router.Post(options.BaseURL+"/api/auth/signup", wrapper.Signup)

	router.Get(options.BaseURL+"/api/health", wrapper.Health)

	router.Get(options.BaseURL+"/api/user", wrapper.GetCurrentUser)

}

type LoginRequestObject struct {
	Body *LoginJSONRequestBody
}

type LoginResponseObject interface {
	VisitLoginResponse(ctx *fiber.Ctx) error
}

type Login200JSONResponse LoginResponse

func (response Login200JSONResponse) VisitLoginResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type LogoutRequestObject struct {
}

type LogoutResponseObject interface {
	VisitLogoutResponse(ctx *fiber.Ctx) error
}

type Logout200Response struct {
}

func (response Logout200Response) VisitLogoutResponse(ctx *fiber.Ctx) error {
	ctx.Status(200)
	return nil
}

type SignupRequestObject struct {
	Body *SignupJSONRequestBody
}

type SignupResponseObject interface {
	VisitSignupResponse(ctx *fiber.Ctx) error
}

type Signup200JSONResponse SignupResponse

func (response Signup200JSONResponse) VisitSignupResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type HealthRequestObject struct {
}

type HealthResponseObject interface {
	VisitHealthResponse(ctx *fiber.Ctx) error
}

type Health200JSONResponse HealthCheckResponse

func (response Health200JSONResponse) VisitHealthResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

type GetCurrentUserRequestObject struct {
}

type GetCurrentUserResponseObject interface {
	VisitGetCurrentUserResponse(ctx *fiber.Ctx) error
}

type GetCurrentUser200JSONResponse GetUserResponse

func (response GetCurrentUser200JSONResponse) VisitGetCurrentUserResponse(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(200)

	return ctx.JSON(&response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (POST /api/auth/login)
	Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error)

	// (GET /api/auth/logout)
	Logout(ctx context.Context, request LogoutRequestObject) (LogoutResponseObject, error)

	// (POST /api/auth/signup)
	Signup(ctx context.Context, request SignupRequestObject) (SignupResponseObject, error)

	// (GET /api/health)
	Health(ctx context.Context, request HealthRequestObject) (HealthResponseObject, error)

	// (GET /api/user)
	GetCurrentUser(ctx context.Context, request GetCurrentUserRequestObject) (GetCurrentUserResponseObject, error)
}

type StrictHandlerFunc func(ctx *fiber.Ctx, args interface{}) (interface{}, error)

type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// Login operation middleware
func (sh *strictHandler) Login(ctx *fiber.Ctx) error {
	var request LoginRequestObject

	var body LoginJSONRequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.Login(ctx.UserContext(), request.(LoginRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Login")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(LoginResponseObject); ok {
		if err := validResponse.VisitLoginResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Logout operation middleware
func (sh *strictHandler) Logout(ctx *fiber.Ctx) error {
	var request LogoutRequestObject

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.Logout(ctx.UserContext(), request.(LogoutRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Logout")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(LogoutResponseObject); ok {
		if err := validResponse.VisitLogoutResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Signup operation middleware
func (sh *strictHandler) Signup(ctx *fiber.Ctx) error {
	var request SignupRequestObject

	var body SignupJSONRequestBody
	if err := ctx.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	request.Body = &body

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.Signup(ctx.UserContext(), request.(SignupRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Signup")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(SignupResponseObject); ok {
		if err := validResponse.VisitSignupResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Health operation middleware
func (sh *strictHandler) Health(ctx *fiber.Ctx) error {
	var request HealthRequestObject

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.Health(ctx.UserContext(), request.(HealthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Health")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(HealthResponseObject); ok {
		if err := validResponse.VisitHealthResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// GetCurrentUser operation middleware
func (sh *strictHandler) GetCurrentUser(ctx *fiber.Ctx) error {
	var request GetCurrentUserRequestObject

	handler := func(ctx *fiber.Ctx, request interface{}) (interface{}, error) {
		return sh.ssi.GetCurrentUser(ctx.UserContext(), request.(GetCurrentUserRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetCurrentUser")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if validResponse, ok := response.(GetCurrentUserResponseObject); ok {
		if err := validResponse.VisitGetCurrentUserResponse(ctx); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}
