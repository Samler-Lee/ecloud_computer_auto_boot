package conf

var Secret = &secret{
	Type: "public",
}

var Server = &server{
	Debug:    false,
	LogLevel: "info",
}

var Cron = &cron{
	Duration: 60,
	Machines: []string{},
}
