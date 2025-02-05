package rest

type QSchemaIn struct {
	Message *string `json:"message,omitempty"`
}

type QSchemaOut struct {
	Message string `json:"message"`
}

type QErrorOut struct {
	Reason string `json:"errReason"`
}
