package menu

import (
	"github.com/YuZongYangHi/spiderX/spider-api/internal/app/apiserver/controllers/base"
	"github.com/YuZongYangHi/spiderX/spider-api/internal/app/apiserver/models"
	"github.com/YuZongYangHi/spiderX/spider-api/internal/app/apiserver/rbac"
	"github.com/YuZongYangHi/spiderX/spider-api/internal/app/apiserver/services"
	"github.com/labstack/echo/v4"
)

func UserPermissionsList(c echo.Context) error {
	cc := c.(*base.Context)
	menuId := cc.ParseInt("id")
	if menuId == 0 {
		return base.BadRequestResponse(c, "")
	}

	userSerializerManager := NewSerializerUserPermissions()
	userSerializerManager.Filter["menu_id"] = menuId

	serializer := base.NewSerializersManager(c, userSerializerManager)
	result, err := serializer.QuerySet()
	if err != nil {
		return base.BadRequestResponse(c, "")
	}

	return base.Response(c, base.HTTP200Code, result)
}

func UserPermissionsCreate(c echo.Context) error {
	menuManager := rbac.NewMenuGrantPermissionsManager(c, services.NewMenuUserGrant())
	response, err := menuManager.Create()
	if err != nil {
		return base.ServerInternalErrorResponse(c, err.Error())
	}
	return base.SuccessResponse(c, response)
}

func UserPermissionsDelete(c echo.Context) error {
	menuManager := rbac.NewMenuGrantPermissionsManager(c, services.NewMenuUserGrant())
	err := menuManager.Delete()
	if err != nil {
		return base.ServerInternalErrorResponse(c, err.Error())
	}
	return base.SuccessNoContentResponse(c)
}

func UserPermissionsUpdate(c echo.Context) error {
	menuManager := rbac.NewMenuGrantPermissionsManager(c, services.NewMenuUserGrant())
	params, err := menuManager.GetParams()
	if err != nil {
		return base.BadRequestResponse(c, "")
	}
	var payload models.PermissionsAction
	valid := base.NewValidator(c)
	if err = valid.IsValid(&payload); err != nil {
		return base.ErrorResponse(c, 400, err.Error())
	}

	obj, err := models.RBACMenuPermissionsUserModel.GetByMenuIdAndId(params.MenuId, params.PermissionsId)
	if err != nil || obj.Id == 0 {
		return base.ErrorResponse(c, 404, err.Error())
	}

	err = models.RBACMenuPermissionsUserModel.UpdateById(obj.Id, &payload)
	if err != nil {
		return base.ServerInternalErrorResponse(c, err.Error())
	}
	return base.SuccessResponse(c, payload)
}
