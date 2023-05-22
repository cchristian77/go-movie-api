package response

/*
	Response structure when displaying a success message.
*/

type Success struct {
	Message string `json:"message"`
	Status  int    `json:"status,omitempty"`
}

var (
	DeleteSuccess = Success{Message: "Data successfully deleted!"}
	UpdateSuccess = Success{Message: "Data successfully updated!"}
)
