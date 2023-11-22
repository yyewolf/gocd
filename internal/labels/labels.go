package labels

type GoCDLabels struct {
	Enable bool   `json:"gocd.enable"`
	Token  string `json:"gocd.token"`
	Repo   string `json:"gocd.repo"`
}

func MapToGoCDLabels(labels map[string]string) GoCDLabels {
	return GoCDLabels{
		Enable: labels["gocd.enable"] == "true",
		Token:  labels["gocd.token"],
		Repo:   labels["gocd.repo"],
	}
}
