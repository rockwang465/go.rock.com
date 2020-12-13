package conf

type CmdConfig struct {
	ServerPort int64  `json:"server_port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	ConfigPath string `json:"config_path"`

	KubeConfig  string `json:"kube_config"`
	TillerPort  int64  `json:"tiller_port"`
	TokenExpire int64  `json:"token_expire"`

	LogDir   string `json:"log_dir"`
	LogLevel string `json:"log_level"`

	DbDriver   string `json:"db_driver"`
	DbPort     int64  `json:"db_port"`
	DbUser     string `json:"db_user"`
	DbHost     string `json:"db_host"`
	DbPassword string `json:"db_password"`
}

func GetCmdCfg() *CmdConfig {
	return &CmdConfig{}
}
