package httpserver

var (
	defaultOption = option{
		ReadTimeout:    10,
		WriteTimeout:   10,
		MaxHeaderBytes: 20,
		Port:           ":9050",
		handles:        make([]Handler, 0),
	}
)
