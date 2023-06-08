package configs

var Env Config

type Config struct {
	App struct {
		Name    string `koanf:"name"`
		Env     string `koanf:"env"`
		Version string `koanf:"version"`
		Port    int32  `koanf:"port"`
	} `koanf:"app"`
	Database struct {
		Host     string `koanf:"host"`
		Port     int32  `koanf:"port"`
		User     string `koanf:"user"`
		Password string `koanf:"password"`
		DBName   string `koanf:"db_name"`
	} `koanf:"database"`
	Context struct {
		Timeout string `koanf:"timeout"`
	} `koanf:"context"`
	Auth struct {
		AccessTokenExpiration  string `koanf:"access_token_expiration"`
		RefreshTokenExpiration string `koanf:"refresh_token_expiration"`
	} `koanf:"auth"`
	JWTKey string `koanf:"jwt_key"`
}
