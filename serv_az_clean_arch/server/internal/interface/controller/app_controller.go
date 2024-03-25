package controller

type AppController struct {
	HandlerController interface{ HandleController }
}
