package labels

type GoCDLabels struct {
	Enable      bool   `json:"gocd.enable"`
	Token       string `json:"gocd.token"`
	Repo        string `json:"gocd.repo"`
	GitlabToken string `json:"gocd.gitlab.token"`
}

func MapToGoCDLabels(labels map[string]string) GoCDLabels {
	return GoCDLabels{
		Enable:      labels["gocd.enable"] == "true",
		Token:       labels["gocd.token"],
		Repo:        labels["gocd.repo"],
		GitlabToken: labels["gocd.gitlab.token"],
	}
}
