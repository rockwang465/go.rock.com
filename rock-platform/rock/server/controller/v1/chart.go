package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/clients/museum"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type ChartVersion struct {
	Name        string       `json:"name" binding:"required" example:"mysql"`
	Version     string       `json:"version" binding:"required" example:"5.7.28-master-54b0c26"`
	Description string       `json:"description" example:"Fast, reliable, scalable, and easy to use open-source relational database system."`
	Keywords    []string     `json:"keywords" example:"mysql,database,sql"`
	Maintainers []Maintainer `json:"maintainers"`
	ApiVersion  string       `json:"apiVersion" example:"v1"`
	AppVersion  string       `json:"appVersion" example:"5.7.28"`
	Urls        []string     `json:"urls" binding:"required" example:"charts/mysql-5.7.28-master-54b0c26.tgz"`
	Created     string       `json:"created" binding:"required" example:"2020-10-23T08:27:01.937112605Z"`
	Digest      string       `json:"digest" binding:"required" example:"76a25ee9205f22c1c922a54a88a161472c1966a54e9d483f16e960449a134ef3"`
}

type Maintainer struct {
	Name  string `json:"name" example:"someone"`
	Email string `json:"email" example:"someone@email.com"`
}

type ChartVersionList []*ChartVersion

type ChartDetail struct {
	Name     string          `json:"name" binding:"required" example:"mysql"`
	Versions []*ChartVersion `json:"version" binding:"required"`
}

type ChartReq struct {
	Name string `json:"name" uri:"name" binding:"required" example:"mysql"`
}

type ChartNameVersionReq struct {
	Name    string `json:"name" uri:"name" binding:"required" example:"mysql"`
	Version string `json:"version" uri:"version" binding:"required" example:"5.7.28-master-54b0c26"`
}

// @Summary Get all chart list
// @Description Api for get all chart list
// @Tags CHART
// @Accept json
// @Produce json
// @Success 200 {array} v1.ChartDetail "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/charts [get]
func (c *Controller) GetAllCharts(ctx *gin.Context) {
	client := museum.GetMuseumClient()
	chartsMapper, err := client.Charts()
	if err != nil {
		panic(err)
	}

	charts := RefactorChartStruct(chartsMapper)

	resp := []*ChartDetail{}
	err = utils.MarshalResponse(charts, &resp)
	if err != nil {
		panic(err)
	}
	c.Logger.Infof("Get all charts list successful, list length %v", len(resp))
	ctx.JSON(http.StatusOK, resp)
}

// format ChartMapper to ChartDetail
func RefactorChartStruct(chartMapper *museum.ChartMapper) []*museum.ChartDetail {
	chartsDetail := make([]*museum.ChartDetail, 0)
	for name, chartVersion := range *chartMapper {
		chartDetail := &museum.ChartDetail{
			Name:     name,
			Versions: chartVersion,
		}
		chartsDetail = append(chartsDetail, chartDetail)
	}
	return chartsDetail
}

// @Summary Get named chart version list
// @Description Api for get named chart version list
// @Tags CHART
// @Accept json
// @Produce json
// @Param name path string true "Chart name"
// @Success 200 {array} v1.ChartVersion "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/charts/{name} [get]
func (c *Controller) GetNamedChartVersions(ctx *gin.Context) {
	var uriReq ChartReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	client := museum.GetMuseumClient()
	chartVersionList, err := client.Versions(uriReq.Name)
	if err != nil {
		panic(err)
	}

	resp := ChartVersionList{}
	if err := utils.MarshalResponse(chartVersionList, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Get %v chart's version list", uriReq.Name)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get named chart specific version info
// @Description Api for get named chart specific version info
// @Tags CHART
// @Accept json
// @Produce json
// @Param name path string true "Chart name"
// @Param version path string true "Chart version"
// @Success 200 {object} v1.ChartVersion "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/charts/{name}/versions/{version} [get]
func (c *Controller) GetNamedChartVersion(ctx *gin.Context) {
	var uriReq ChartNameVersionReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	client := museum.GetMuseumClient()
	chartVersion, err := client.Version(uriReq.Name, uriReq.Version)
	if err != nil {
		panic(err)
	}

	resp := ChartVersion{}
	if err := utils.MarshalResponse(chartVersion, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Get chart(%v) version(%v) info", uriReq.Name, uriReq.Version)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Delete named chart specific version
// @Description Api for delete named chart specific version
// @Tags CHART
// @Accept json
// @Produce json
// @Param name path string true "Chart name"
// @Param version path string true "Chart version"
// @Success 200 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/charts/{name}/versions/{version} [delete]
func (c *Controller) DeleteNamedChartVersion(ctx *gin.Context) {
	var uriReq ChartNameVersionReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	client := museum.GetMuseumClient()
	if err := client.DeleteVersion(uriReq.Name, uriReq.Version); err != nil {
		panic(err)
	}
	c.Logger.Infof("Delete chart:%v version:%v", uriReq.Name, uriReq.Version)
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Get specific app's all charts version list
// @Description Api to get an app's all charts version list
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Success 200 {array} v1.ChartVersion "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id}/charts [get]
func (c *Controller) GetAppChartVersions(ctx *gin.Context) {
	var uriReq IdReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	app, err := api.GetAppById(uriReq.Id)
	if err != nil {
		panic(err)
	}
	client := museum.GetMuseumClient()
	chartVersionList, err := client.Versions(app.Name)
	if err != nil {
		panic(err)
	}

	resp := ChartVersionList{}
	if err := utils.MarshalResponse(chartVersionList, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Get %v chart's version list by app id:%v", app.Name, uriReq.Id)
	ctx.JSON(http.StatusOK, resp)
}
