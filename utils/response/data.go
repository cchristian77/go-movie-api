package response

/*
	Response structure when displaying data
*/

type Result struct {
	Meta any `json:"meta,omitempty"`
	Data any `json:"data"`
}
