//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.cfg.yaml ../spec.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../spec.yaml

package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Api struct{}

func (a *Api) FindUserById(ctx echo.Context, userId UserId) error {
	user := User{Id: int(userId), Firstname: nil, Surname: "Doe"}
	return ctx.JSON(http.StatusOK, user)
}
