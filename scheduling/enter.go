package scheduling

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"time"
)

func CronInit() {
	timezone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logrus.Error(err.Error())
	}
	Cron := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	_, _ = Cron.AddFunc("0 */60 * * * *", RunCftTask)
	Cron.Start()
}
