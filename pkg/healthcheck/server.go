package healthcheck

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"port-checker/pkg/logging"
)

type paramsType struct {
}

// Serve healthcheck HTTP requests and listens on
// localhost:9999 only.
func Serve() {
	params := paramsType{}
	localRouter := httprouter.New()
	localRouter.GET("/", params.get)
	logging.Info("Private server listening on 127.0.0.1:9999")
	err := http.ListenAndServe("127.0.0.1:9999", localRouter)
	if err != nil {
		logging.Fatal("%s", err)
	}
}

func (params *paramsType) get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO ping main server
	w.WriteHeader(http.StatusOK)
}
