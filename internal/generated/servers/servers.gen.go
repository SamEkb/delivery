// Package servers provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package servers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Courier defines model for Courier.
type Courier struct {
	// Id Идентификатор
	Id       openapi_types.UUID `json:"id"`
	Location Location           `json:"location"`

	// Name Имя
	Name string `json:"name"`
}

// Error defines model for Error.
type Error struct {
	// Code Код ошибки
	Code int32 `json:"code"`

	// Message Текст ошибки
	Message string `json:"message"`
}

// Location defines model for Location.
type Location struct {
	// X X
	X int `json:"x"`

	// Y Y
	Y int `json:"y"`
}

// NewCourier defines model for NewCourier.
type NewCourier struct {
	// Name Имя
	Name string `json:"name"`

	// Speed Скорость
	Speed int `json:"speed"`
}

// Order defines model for Order.
type Order struct {
	// Id Идентификатор
	Id       openapi_types.UUID `json:"id"`
	Location Location           `json:"location"`
}

// CreateCourierJSONRequestBody defines body for CreateCourier for application/json ContentType.
type CreateCourierJSONRequestBody = NewCourier

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Получить всех курьеров
	// (GET /api/v1/couriers)
	GetCouriers(ctx echo.Context) error
	// Добавить курьера
	// (POST /api/v1/couriers)
	CreateCourier(ctx echo.Context) error
	// Создать заказ
	// (POST /api/v1/orders)
	CreateOrder(ctx echo.Context) error
	// Получить все незавершенные заказы
	// (GET /api/v1/orders/active)
	GetOrders(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetCouriers converts echo context to params.
func (w *ServerInterfaceWrapper) GetCouriers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCouriers(ctx)
	return err
}

// CreateCourier converts echo context to params.
func (w *ServerInterfaceWrapper) CreateCourier(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateCourier(ctx)
	return err
}

// CreateOrder converts echo context to params.
func (w *ServerInterfaceWrapper) CreateOrder(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateOrder(ctx)
	return err
}

// GetOrders converts echo context to params.
func (w *ServerInterfaceWrapper) GetOrders(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetOrders(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/api/v1/couriers", wrapper.GetCouriers)
	router.POST(baseURL+"/api/v1/couriers", wrapper.CreateCourier)
	router.POST(baseURL+"/api/v1/orders", wrapper.CreateOrder)
	router.GET(baseURL+"/api/v1/orders/active", wrapper.GetOrders)

}

type GetCouriersRequestObject struct {
}

type GetCouriersResponseObject interface {
	VisitGetCouriersResponse(w http.ResponseWriter) error
}

type GetCouriers200JSONResponse []Courier

func (response GetCouriers200JSONResponse) VisitGetCouriersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetCouriersdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response GetCouriersdefaultJSONResponse) VisitGetCouriersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type CreateCourierRequestObject struct {
	Body *CreateCourierJSONRequestBody
}

type CreateCourierResponseObject interface {
	VisitCreateCourierResponse(w http.ResponseWriter) error
}

type CreateCourier201Response struct {
}

func (response CreateCourier201Response) VisitCreateCourierResponse(w http.ResponseWriter) error {
	w.WriteHeader(201)
	return nil
}

type CreateCourier400JSONResponse Error

func (response CreateCourier400JSONResponse) VisitCreateCourierResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type CreateCourier409JSONResponse Error

func (response CreateCourier409JSONResponse) VisitCreateCourierResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type CreateCourierdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response CreateCourierdefaultJSONResponse) VisitCreateCourierResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type CreateOrderRequestObject struct {
}

type CreateOrderResponseObject interface {
	VisitCreateOrderResponse(w http.ResponseWriter) error
}

type CreateOrder201Response struct {
}

func (response CreateOrder201Response) VisitCreateOrderResponse(w http.ResponseWriter) error {
	w.WriteHeader(201)
	return nil
}

type CreateOrderdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response CreateOrderdefaultJSONResponse) VisitCreateOrderResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetOrdersRequestObject struct {
}

type GetOrdersResponseObject interface {
	VisitGetOrdersResponse(w http.ResponseWriter) error
}

type GetOrders200JSONResponse []Order

func (response GetOrders200JSONResponse) VisitGetOrdersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetOrdersdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response GetOrdersdefaultJSONResponse) VisitGetOrdersResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Получить всех курьеров
	// (GET /api/v1/couriers)
	GetCouriers(ctx context.Context, request GetCouriersRequestObject) (GetCouriersResponseObject, error)
	// Добавить курьера
	// (POST /api/v1/couriers)
	CreateCourier(ctx context.Context, request CreateCourierRequestObject) (CreateCourierResponseObject, error)
	// Создать заказ
	// (POST /api/v1/orders)
	CreateOrder(ctx context.Context, request CreateOrderRequestObject) (CreateOrderResponseObject, error)
	// Получить все незавершенные заказы
	// (GET /api/v1/orders/active)
	GetOrders(ctx context.Context, request GetOrdersRequestObject) (GetOrdersResponseObject, error)
}

type StrictHandlerFunc = strictecho.StrictEchoHandlerFunc
type StrictMiddlewareFunc = strictecho.StrictEchoMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// GetCouriers operation middleware
func (sh *strictHandler) GetCouriers(ctx echo.Context) error {
	var request GetCouriersRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetCouriers(ctx.Request().Context(), request.(GetCouriersRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetCouriers")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetCouriersResponseObject); ok {
		return validResponse.VisitGetCouriersResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// CreateCourier operation middleware
func (sh *strictHandler) CreateCourier(ctx echo.Context) error {
	var request CreateCourierRequestObject

	var body CreateCourierJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.CreateCourier(ctx.Request().Context(), request.(CreateCourierRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateCourier")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(CreateCourierResponseObject); ok {
		return validResponse.VisitCreateCourierResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// CreateOrder operation middleware
func (sh *strictHandler) CreateOrder(ctx echo.Context) error {
	var request CreateOrderRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.CreateOrder(ctx.Request().Context(), request.(CreateOrderRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateOrder")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(CreateOrderResponseObject); ok {
		return validResponse.VisitCreateOrderResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// GetOrders operation middleware
func (sh *strictHandler) GetOrders(ctx echo.Context) error {
	var request GetOrdersRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetOrders(ctx.Request().Context(), request.(GetOrdersRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetOrders")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetOrdersResponseObject); ok {
		return validResponse.VisitGetOrdersResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9xVT28bRRT/KqsHx1Vtt1zYIwEhpIgeuIAQh8V+caby7iyz47RWtFLsQBspUXPhgHoo",
	"KnyBjcnKi0Pcr/DeN0IzEyduPHFcEVURF8u7O/Pm9++92YW2TDKZYqpziHYhb29jEtu/G7KvBCrzN1My",
	"Q6UF2g+iY347mLeVyLSQKURAv9EpVXTOI6r5Z6ppSiWPaMZ7EMKWVEmsIYJ+X3QgBD3IECLItRJpF4oQ",
	"erIdu0K78LHCLYjgo8YVsMYFqsbmfF0RQhon6MXxDx8vn1GEoPCnvlDYgeh7sDBshYXDf7jcJX98gm1t",
	"TvlCKemRoC07vsNf0YxOA5rxAdV0QlOqF9mLVD96eAVNpBq7qMwpCeZ53PVV/IMqmvKQR9erruZn8V3V",
	"9THbXND8XXLPlnF8a4qJVCT9BKKmj8JgedN3t2y6hvkZmCo+qF/j0xvDeEsMEpFuYtrV2xC1PMHLM0Rf",
	"mt/Q1GSXZkZ6Plok0rqVyEWuXG0fn8eqc1/7ytcnKxrErBfplvQAf80jGlPFL6ikyuR3QmXA+/zCPU15",
	"n/f4iCqjMo3DwPDkIb01n+2iParNHn5ONb80n60ZVNKYZjQN330z5X0jgNA9A++bp3G3iyr4HHtiB9UA",
	"QthBlTtkrQfNB02jjswwjTMBETyyr0LIYr1trWjEmWjstBptFzv7rovaQ/N3mtHEQjrjY0ftrX0wTGsT",
	"noDGPKSKf1kiDRaDsuJ+1YEIvkS9MT/RGJFnMs1dOB42m27wpBpTCyTOsp5wzjSe5M5kZ6gNk8Ykv833",
	"eVcVl8bGSsUD5+s1on9emHNA53xIf5t55AwegV28Ffd7+r0grkLm5q4Px+vLMVjauOb9JInVYO7FesIX",
	"IWQyX9PPU5rRiU3ZRdnFauWSiRsKY41zaV0/Ya4/k53BncmzMBF9Gr26AgjFUpBaHtqr3f3kPcP3n50N",
	"aEwlnVFNp24CUO1wfPrBcfCha2g6N3OYaj4O6MSOpnMzsAI6oxn9ZSdzfW864dfVkTWr5yNOmsvIJmPt",
	"juChfWWsscUnVBoQNAl4GPBzquiMj/hlwCOjkL26bNtR6QS8oWXctXgHcb0XFry5QSOP+I24rcUO3sEl",
	"E9iUTqzzFe/xgc2s0ahagMCHvpvnsQvCh7h3nNP/61tnbSeKoij+DQAA//8JDQamhg0AAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
