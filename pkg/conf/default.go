package conf

var Secret = &secret{}

var Server = &server{
	Debug:    false,
	LogLevel: "info",
}

var Cron = &cron{
	Duration:    60,
	MachineList: []string{},
}
