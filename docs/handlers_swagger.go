package docs


import "github.com/PhilLar/Images-back/handlers"

// swagger:route POST /files files-tag idOfFilesEndpoint
// /files does some amazing stuff.
// responses:
//   200: filesResponse
//   400: HTTPError

// This text will appear as description of your response body.
// swagger:response filesResponse
type filesResponseWrapper struct {
	// in:body
	Image struct {
		Body handlers.ImageFile
	}	`json:"image"`
}
