package docs

import (
	"github.com/PhilLar/Images-back/handlers"
)

// swagger:route POST /files files idOfFilesEndpoint
// /files upload your image to the server.
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
// responses:
//   200: filesResponse
//   400: description: invalid type of file (image)

// This text will appear as description of your response body.
// swagger:response filesResponse
type filesResponseWrapper struct {
	// in:body
	Body struct {
		handlers.ImageFile
	}
}

// swagger:parameters idOfFilesEndpoint
type filesParamsWrapper struct {
	// ImageFile object that needed to be added to the db.
	// in: body
	// Required: true
	Body struct {
		Title		string	`json:"title"`
		FileName	string	`json:"file_name"`
	}
}

// swagger:route GET /images files idOfImagesEndpoint
// /images endpoint lists all images contained in database.
//     Produces:
//     - application/json
// responses:
//	 200: imagesResponse
//   400: description: Bad Request


// This text will appear as description of your response body.
// swagger:response imagesResponse
type imagesResponseWrapper struct {
	// in:body
	Body []*handlers.ImageFile
}

// swagger:route DELETE /images/{ID} files idOfEndpoint
// /images/{ID} endpoint deletes an image by ID.
//     Produces:
//     - application/json
// responses:
//	 200: imagesResponse
//   400: description: image with such ID not found in '/files' directory
//   400: description: ID must be integer (BIGSERIAL)

