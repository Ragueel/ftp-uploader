package config

type AppAuthConfig struct {
	Username string
	Password string
	Host     string
	Name     string
}

type UploadConfig struct {
	AuthConfig     AppAuthConfig
	LocalRootPath  string
	UploadRootPath string
	IgnorePaths    *[]string
}

func NewAuthConfig() AppAuthConfig {
	return AppAuthConfig{}
}

func NewUploadConfig() UploadConfig {
	return UploadConfig{}
}
