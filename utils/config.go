package utils

type Config struct {
	GiteeToken        string `json:"gitee_token"          required:"true"`
	OpenEulerToken    string `json:"openeuler_token"      required:"true"`
	OpenEulerId       string `json:"openeuler_id"         required:"true"`
	EulerUserSigUrl   string `json:"euler_user_sig_url"   required:"true"`
	GiteeUserReposUrl string `json:"gitee_user_repos_url"   required:"true"`
	GiteeGetUserIdUrl string `json:"gitee_get_user_id_url"   required:"true"`
	GiteeGetPullsUrl  string `json:"gitee_get_pulls_url"   required:"true"`
}

var config Config

func Init(cfg *Config) {
	config = *cfg
}
