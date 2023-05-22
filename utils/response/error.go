package response

/*
	Response structure when displaying an error.
*/

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status,omitempty"`
	Error   error  `json:"error,omitempty"`
}
