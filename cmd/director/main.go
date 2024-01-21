package main

import (
	"github.com/johnnyluo/tss-director/handler"
	"github.com/johnnyluo/tss-director/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("10M")) // set maximum allowed size for a request body to 10M
	store := storage.NewInMemoryStorage()
	s := handler.NewServer(store)
	e.GET("/:sessionID", s.StartSession)
	e.DELETE("/:sessionID", s.EndSession)
	e.GET(("/message/:sessionID/:participantID"), s.GetMessage)
	e.POST(("/message/:sessionID"), s.PostMessage)
	e.Logger.Fatal(e.Start(":8080"))
}
