//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../spec.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../spec.yaml

package api

import (
  "encoding/json"
	"net/http"
)


type ApiImpl struct{}

func (a *ApiImpl) FindUserById(w http.ResponseWriter, r *http.Request, userId UserId) {
	user := User{Id: int(userId), Firstname: nil, Surname: "Doe"}
  w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func SetupHandler() {
    api := ApiImpl{}
    http.Handle("/", Handler(&api))
}
