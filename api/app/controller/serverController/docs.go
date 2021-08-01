//lint:file-ignore U1000 Ignore all unused code

package serverController

// swagger:route GET /ping server ping
//
// ping
//
// Check if the server is running.
//
// responses:
//   200: pingResponse

// Server is running!
// swagger:response pingResponse
type pingResponseWrapper struct {
	// in:body
	Body PingResponse
}

// swagger:route GET /version server version
//
// version
//
// Returns the version of the api.
//
// responses:
//   200: versionResponse

// swagger:response versionResponse
type versionResponseWrapper struct {
	// in:body
	Body VersionResponse
}
