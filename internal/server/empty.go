package server

import "net/http"

type emptyOKResp struct{}

func (e *emptyOKResp) writeJSON(w http.ResponseWriter) error {
	return writeJSONResponse(w, http.StatusOK, e)
}

func getEmptyOKResponse() response {
	return &emptyOKResp{}
}
