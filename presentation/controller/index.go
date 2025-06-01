package controller

var Controllers = []Controller{}

func Add(controller Controller) {
	Controllers = append(Controllers, controller)
}

func SetupRoutes() {
	for _, controller := range Controllers {
		controller.SetupRoutes()
	}
}