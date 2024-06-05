package config

type CftConfig struct {
	Root         string `yaml:"root"`
	Execute      string `yaml:"execute"`
	Result       string `yaml:"result"`
	Key          string `yaml:"key"`
	AnalysisName string `yaml:"analysis_name"`
	AnalysisType string `yaml:"analysis_type"`
	Email        string `yaml:"email"`
	Url          string `yaml:"url"`
}
