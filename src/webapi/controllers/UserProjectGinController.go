package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rol/app/services"
	"rol/dtos"
	"rol/webapi"
)

//UserProjectController user project API controller
type UserProjectController struct {
	service *services.UserProjectService
	logger  *logrus.Logger
}

//NewUserProjectController user project controller constructor. Parameters pass through DI
//
//Params:
//	projectService - user project service
//	log - logrus logger
//Return:
//	*UserProjectController - instance of user project controller
func NewUserProjectController(projectService *services.UserProjectService, log *logrus.Logger) *UserProjectController {
	return &UserProjectController{
		service: projectService,
		logger:  log,
	}
}

//RegisterUserProjectController registers controller for getting user projects via api
func RegisterUserProjectController(controller *UserProjectController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")

	groupRoute.GET("/project/", controller.GetList)
	groupRoute.GET("/project/:id", controller.GetByID)
	groupRoute.POST("/project/", controller.Create)
	groupRoute.DELETE("/project/:id", controller.Delete)
}

//GetList get list of user projects
//
//Params:
//	ctx - gin context
//
// @Summary Get list of user projects
// @version	1.0
// @Tags	project
// @Accept	json
// @Produce	json
// @Success	200		{object}	[]dtos.ProjectDto
// @Failure	500		"Internal Server Error"
// @router	/project/	[get]
func (p *UserProjectController) GetList(ctx *gin.Context) {
	projList, err := p.service.GetList(ctx, "", "", "", 1, 10)
	handleWithData(ctx, err, projList)
}

//GetByID get user project by id
//
//Params:
//	ctx - gin context
//
// @Summary	Gets user project by id
// @version	1.0
// @Tags	project
// @Accept	json
// @Produce	json
// @param	id	path		string	true	"User project ID"
// @Success	200		{object}	dtos.ProjectDto
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router	/project/{id}	[get]
func (p *UserProjectController) GetByID(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	vlan, err := p.service.GetByID(ctx, id)
	handleWithData(ctx, err, vlan)
}

//Create new user project
//
//Params:
//	ctx - gin context
//
// @Summary	Create new user project
// @version	1.0
// @Tags	project
// @Accept	json
// @Produce	json
// @Param	request	body		dtos.ProjectCreateDto	true	"User project fields"
// @Success	200		{object}	dtos.ProjectDto
// @Failure	400		{object}	dtos.ValidationErrorDto
// @Failure	500		"Internal Server Error"
// @router	/project/	[post]
func (p *UserProjectController) Create(ctx *gin.Context) {
	reqDto, err := getRequestDtoAndRestoreBody[dtos.ProjectCreateDto](ctx)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}

	vlanDto, err := p.service.Create(ctx, reqDto)
	handleWithData(ctx, err, vlanDto)
}

//Delete user project
//
//Params:
//	ctx - gin context
//
// @Summary	Delete user project by id
// @version	1.0
// @Tags	project
// @Accept	json
// @Produce	json
// @param	id	path		string	true	"User project ID"
// @Success	204		"OK, but No Content"
// @Failure	404		"Not Found"
// @Failure	500		"Internal Server Error"
// @router	/project/{id}	[delete]
func (p *UserProjectController) Delete(ctx *gin.Context) {
	strID := ctx.Param("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		abortWithStatusByErrorType(ctx, err)
		return
	}
	err = p.service.Delete(ctx, id)
	handle(ctx, err)
}
