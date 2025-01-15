package models

// Redirect represents a redirection entity with a URL and an optional key.
// The URL field specifies the target destination and is required with validation as a valid URL.
// The Key field is optional and can be used to uniquely identify the redirection.
type Redirect struct {
	URL string `json:"url" binding:"required,url"`
	Key string `json:"key,omitempty" binding:"-"`
}
