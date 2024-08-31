package weekly_report_app

import "wbapp/internal/config"

const configPath = "../../configs/config.yaml"

func Run() error {
	_, err := config.ParseConfig(configPath)
	if err != nil {
		return err
	}
	return nil
}
