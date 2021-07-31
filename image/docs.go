package image

//
// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handlers

// Generic error message returned as a string
// swagger:response imageErrorResponse
//lint:ignore U1000 for docs
type imageErrorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// A image
// swagger:response imageResponse
//lint:ignore U1000 for docs
type imageResponseWrapper struct {
	// A image
	// in: body
	Body []byte
}

// swagger:response imageIdResponse
//lint:ignore U1000 for docs
type imageIdResponseWrapper struct {
	// The id of the image for which the operation relates
	// in: body
	Id string `json:"id"`
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
//lint:ignore U1000 for docs
type noContentResponseWrapper struct {
}

// swagger:parameters getImage deleteImage
//lint:ignore U1000 for docs
type imageIdParamsWrapper struct {
	// The id of the image for which the operation relates
	// in: path
	// required: true
	Id string `json:"id"`
}

// swagger:parameters putImage
//lint:ignore U1000 for docs
type imageParamsWrapper struct {
	// in: body
	// required: true
	Body []byte
}
