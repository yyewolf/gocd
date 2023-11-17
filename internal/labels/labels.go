package labels

type GoCDLabels struct {
	Enable bool   `json:"gocd.enable"`
	Token  string `json:"gocd.token"`
}

func MapToGoCDLabels(labels map[string]string) GoCDLabels {
	return GoCDLabels{
		Enable: labels["gocd.enable"] == "true",
		Token:  labels["gocd.token"],
	}
}
