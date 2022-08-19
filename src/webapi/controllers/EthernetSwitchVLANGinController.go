package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"rol/app/services"
	"rol/dtos"
	"rol/webapi"
	"strconv"

	"github.com/sirupsen/logrus"
)

//EthernetSwitchVLANGinController ethernet switch GIN controller constructor
type EthernetSwitchVLANGinController struct {
	service *services.EthernetSwitchVLANService
	logger  *logrus.Logger
}

//NewEthernetSwitchVLANGinController ethernet switch VLAN controller constructor. Parameters pass through DI
//Params
//	service - generic service
//	log - logrus logger
//Return
//	*GinGenericController - instance of generic controller for http logs
func NewEthernetSwitchVLANGinController(service *services.EthernetSwitchVLANService, log *logrus.Logger) *EthernetSwitchVLANGinController {
	ethernetSwitchVLANController := &EthernetSwitchVLANGinController{
		service: service,
		logger:  log,
	}
	return ethernetSwitchVLANController
}

//RegisterEthernetSwitchVLANGinController registers controller for getting ethernet switch VLANs via api
func RegisterEthernetSwitchVLANGinController(controller *EthernetSwitchVLANGinController, server *webapi.GinHTTPServer) {
	groupRoute := server.Engine.Group("/api/v1")
	groupRoute.GET("/ethernet-switch/:id/vlan/", controller.GetList)
	groupRoute.GET("/ethernet-switch/:id/vlan/:vlanID", controller.GetByID)
	groupRoute.POST("/ethernet-switch/:id/vlan/", controller.Create)
	groupRoute.PUT("/ethernet-switch/:id/vlan/", controller.Update)
	groupRoute.DELETE("/ethernet-switch/:id/vlan/", controller.Delete)
}

//GetList get list of switch VLANs with search and pagination
//	Params
//	ctx - gin context
// @Summary Get paginated list of switch VLANs
// @version 1.0
// @Tags ethernet-switch
// @Accept  json
// @Produce json
// @param 	 id 			 query  string  true "Ethernet switch ID"
// @param	 orderBy		 query	string	false	"Order by field"
// @param	 orderDirection	 query	string	false	"'asc' or 'desc' for ascending or descending order"
// @param	 search			 query	string	false	"Searchable value in entity"
// @param	 page			 query	int		false	"Page number"
// @param	 pageSize		 query	int		false	"Number of entities per page"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.PaginatedListDto{items=[]dtos.EthernetSwitchVLANDto}}
// @router /ethernet-switch/{id}/vlans [get]
func (e *EthernetSwitchVLANGinController) GetList(ctx *gin.Context) {
	orderBy := ctx.DefaultQuery("orderBy", "Name")
	orderDirection := ctx.DefaultQuery("orderDirection", "asc")
	search := ctx.DefaultQuery("search", "")
	page := ctx.DefaultQuery("page", "1")
	pageInt64, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		pageInt64 = 1
	}
	pageSize := ctx.DefaultQuery("pageSize", "10")
	pageSizeInt64, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		pageSizeInt64 = 10
	}
	strSwitchID := ctx.Param("id")
	switchID, err := uuid.Parse(strSwitchID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}

	paginatedList, err := e.service.GetVLANs(ctx, switchID, search, orderBy, orderDirection, int(pageInt64), int(pageSizeInt64))
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			e.logger.Errorf("%s : %s", err, controllerErr)
		}
		return
	}
	responseDto := &dtos.ResponseDataDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
		Data: paginatedList,
	}
	ctx.JSON(http.StatusOK, responseDto)
}

//GetByID get switch VLAN by id
//	Params
//	ctx - gin context
// @Summary Get ethernet switch VLAN by id
// @version 1.0
// @Tags ethernet-switch
// @Accept  json
// @Produce  json
// @param	 id			query	string		true	"Ethernet switch ID"
// @param	 vlanID		query	string		true	"Ethernet switch VLAN ID"
// @Success 200 {object} dtos.ResponseDataDto{data=dtos.EthernetSwitchVLANDto}
// @router /ethernet-switch/{id}/vlans/{vlanID} [get]
func (e *EthernetSwitchVLANGinController) GetByID(ctx *gin.Context) {
	strID := ctx.Param("vlanID")
	id, err := uuid.Parse(strID)
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusNotFound, err)
		if controllerErr != nil {
			e.logger.Errorf("%s : %s", err, controllerErr)
		}
	}
	strSwitchID := ctx.Param("id")
	switchID, err := uuid.Parse(strSwitchID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}
	dto, err := e.service.GetVLANByID(ctx, switchID, id)
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			e.logger.Errorf("%s : %s", err, controllerErr)
		}
	}
	if dto == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}
	responseDto := &dtos.ResponseDataDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
		Data: dto,
	}
	ctx.JSON(http.StatusOK, responseDto)
}

//Create new switch VLAN
//	Params
//	ctx - gin context
// @Summary Create new ethernet switch VLAN
// @version 1.0
// @Tags ethernet-switch
// @Accept  json
// @Produce  json
// @param id query string true "Ethernet switch ID"
// @Param request body dtos.EthernetSwitchVLANCreateDto true "Ethernet switch VLAN fields"
// @Success 200 {object} dtos.ResponseDataDto{data=uuid.UUID}
// @router /ethernet-switch/{id}/vlans [post]
func (e *EthernetSwitchVLANGinController) Create(ctx *gin.Context) {
	reqDto := new(dtos.EthernetSwitchVLANCreateDto)
	err := ctx.ShouldBindJSON(&reqDto)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	// Restoring body in gin.Context for logging it later in middleware
	err = RestoreBody(reqDto, ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	strSwitchID := ctx.Param("id")
	switchID, err := uuid.Parse(strSwitchID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}
	id, err := e.service.CreateVLAN(ctx, switchID, *reqDto)
	if err != nil {
		controllerErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if controllerErr != nil {
			e.logger.Errorf("%s : %s", err, controllerErr)
		}
		return
	}
	responseDto := dtos.ResponseDataDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
		Data: id,
	}
	ctx.JSON(http.StatusOK, responseDto)
}

//Update switch VLAN by id
//	Params
//	ctx - gin context
// @Summary Updates ethernet switch VLAN by id
// @version 1.0
// @Tags ethernet-switch
// @Accept  json
// @Produce  json
// @param id query string true "Ethernet switch ID"
// @Param request body dtos.EthernetSwitchVLANUpdateDto true "Ethernet switch fields"
// @Success 200 {object} dtos.ResponseDto
// @router /ethernet-switch/{id}/vlans [put]
func (e *EthernetSwitchVLANGinController) Update(ctx *gin.Context) {
	reqDto := new(dtos.EthernetSwitchVLANUpdateDto)
	err := ctx.ShouldBindJSON(reqDto)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	// Restoring body in gin.Context for logging it later in middleware
	err = RestoreBody(reqDto, ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	strSwitchID := ctx.Param("id")
	switchID, err := uuid.Parse(strSwitchID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}

	strPortID := ctx.Param("portID")
	portID, err := uuid.Parse(strPortID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}

	err = e.service.UpdateVLAN(ctx, switchID, portID, *reqDto)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	responseDto := &dtos.ResponseDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
	}
	ctx.JSON(http.StatusOK, responseDto)
}

//Delete soft deleting switch VLAN in database
//	Params
//	ctx - gin context
// @Summary Delete ethernet switch VLAN by id
// @version 1.0
// @Tags ethernet-switch
// @Accept  json
// @Produce  json
// @param	 id			query	string		true	"Ethernet switch ID"
// @param	 vlanID		query	string		true	"Ethernet switch VLAN ID"
// @Success 200 {object} dtos.ResponseDto
// @router /ethernet-switch/{id}/vlans/{vlanID}  [delete]
func (e *EthernetSwitchVLANGinController) Delete(ctx *gin.Context) {
	strSwitchID := ctx.Param("id")
	switchID, err := uuid.Parse(strSwitchID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}

	strPortID := ctx.Param("portID")
	portID, err := uuid.Parse(strPortID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}

	err = e.service.DeleteVLAN(ctx, switchID, portID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	responseDto := &dtos.ResponseDto{
		Status: dtos.ResponseStatusDto{
			Code:    0,
			Message: "",
		},
	}
	ctx.JSON(http.StatusOK, responseDto)
}
