package config

type AppAuthConfig struct {
	Username string
	Password string
	Host     string
}

type UploadConfig struct {
	AuthConfig     AppAuthConfig
	LocalRootPath  string
	UploadRootPath string
	Name           string
	IgnorePaths    *[]string
}

func NewAuthConfig() AppAuthConfig {
	return AppAuthConfig{}
}

func NewUploadConfig() UploadConfig {
	return UploadConfig{}
}