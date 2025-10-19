package productHTTPHandler

import (
	"dijam-ecommerce/internal/product"
	"dijam-ecommerce/pkg/apierrors"
	"dijam-ecommerce/shared"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ProductHTTPHandler struct {
	ProductService *product.ProductService
	validator      *validator.Validate
}

func NewProductHTTPHandler(productService *product.ProductService, validator *validator.Validate) *ProductHTTPHandler {
	return &ProductHTTPHandler{
		ProductService: productService,
		validator:      validator,
	}
}

func formatValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)

	if vErrs, ok := err.(validator.ValidationErrors); ok {
		for _, vErr := range vErrs {
			field := strings.ToLower(vErr.Field())
			switch vErr.Tag() {
			case "required":
				errorsMap[field] = "this field is required"
			case "gte":
				errorsMap[field] = fmt.Sprintf("this field must be greater than or equal to %s", vErr.Param())
			case "lte":
				errorsMap[field] = fmt.Sprintf("this field must be less than or equal to %s", vErr.Param())
			case "oneof":
				errorsMap[field] = fmt.Sprintf("this field must be one of: %s", vErr.Param())
			default:
				errorsMap[field] = fmt.Sprintf("this field is invalid (failed on tag: %s)", vErr.Tag())
			}
		}
	}
	return errorsMap
}

func (h *ProductHTTPHandler) ListProduct(w http.ResponseWriter, r *http.Request) {

	var input struct {
		// Name          string
		// PriceInPence  int64  `validate:"gte=0"`
		Category      string `validate:"oneof=electronics clothing books home_goods sports toys health automotive grocery computers other ''"`
		StockQuantity int64  `validate:"gte=0"`
		Filters       shared.ProductFilters
	}

	qs := r.URL.Query()
	// input.Name = shared.QueryReadString(qs, "name", "")
	// input.PriceInPence = shared.QueryReadInt(qs, "price_in_pence", 0)
	input.Category = shared.QueryReadString(qs, "category", "")
	input.StockQuantity = shared.QueryReadInt(qs, "stock_quantity", 1)
	input.Filters.Sort = shared.QueryReadString(qs, "sort", "price_in_pence")
	input.Filters.Page = shared.QueryReadInt(qs, "page", 1)
	input.Filters.PageSize = shared.QueryReadInt(qs, "page_size", 20)

	err := h.validator.Struct(input)
	if err != nil {
		validationErrors := formatValidationErrors(err)
		apierrors.ValidationFailed(w, r, validationErrors)
		return
	}

	// products, metadata, err := h.ProductService.List(r.Context(), input.Category, input.Filters)
	_, _, err = h.ProductService.List(r.Context(), input.Category, input.Filters)
	if err != nil {
		apierrors.ServerError(w, r, err)
		return
	}
	// fmt.Fprintf(w, "%+v\n", input)
	// shared.WriteToJSON(w, http.StatusOK, shared.Envelope{"message": "successfully retrieved products", "metadata": metadata, "products": products}, nil)
}
