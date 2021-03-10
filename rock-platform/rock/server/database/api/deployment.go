package api

import (
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

// insert the current deployment info into the database
func CreateDeployment(appId, envId int64, chartName, chartVersion, description, namespace string) (*models.Deployment, error) {
	deployment := models.Deployment{
		Description:  description,
		ChartName:    chartName,
		ChartVersion: chartVersion,
		AppId:        appId,
		EnvId:        envId,
	}

	deployment.Name = utils.GenerateChartName(chartName, namespace)

	// 由于 app_id env_id 在调用此函数前已经做了检测了，所以这里就不再做检测了。
	db := database.GetDBEngine()
	if err := db.Create(&deployment).Error; err != nil {
		return nil, err
	}
	return &deployment, nil
}
