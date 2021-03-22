package drone_api

import (
	// "github.com/drone/drone-go/drone" // 之前用的是这个项目的v0.8.4的tag，并修改了源码
	"github.com/rockwang465/drone/drone-go/drone" // 所以不用被改动的源码，而是直接拿改好的源码放在自己的项目下
	"go.rock.com/rock-platform/rock/server/conf"
	"go.rock.com/rock-platform/rock/server/utils"
	"golang.org/x/oauth2"
	"net/http"
)

// use drone token, generate a drone client
func getClient(token string) (drone.Client, error) {
	cfg := conf.GetConfig()
	droneAddr := cfg.Viper.GetString("drone.addr")
	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: token,
		},
	)

	if token == "" { // 如果没有token(drone token),说明此用户还没有关联过access token;因为只有在关联access token时才会生成drone token.
		err := utils.NewRockError(412, 41200001, "Please set gitlab access token first")
		return drone.NewClient(droneAddr, auther), err
	}

	return drone.NewClient(droneAddr, auther), nil
}

func getRawClient() drone.Client {
	cfg := conf.GetConfig()
	droneAddr := cfg.Viper.GetString("drone.addr")
	client := &http.Client{}
	return drone.NewClient(droneAddr, client)
}
