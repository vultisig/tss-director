package handler

import (
	"net/http"
	"strings"

	"github.com/johnnyluo/tss-director/model"
	"github.com/johnnyluo/tss-director/storage"
	"github.com/labstack/echo/v4"
)

type Server struct {
	s storage.Storage
}

// NewServer returns a new server.
func NewServer(s storage.Storage) *Server {
	return &Server{s: s}
}

// StartSession is to start a new session that will be used to send and receive messages.
func (s *Server) StartSession(c echo.Context) error {
	sessionID := c.Param("sessionID")
	if strings.Trim(sessionID, " ") == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	if _, err := s.s.GetSession(sessionID); err != storage.ErrNotFound {
		return c.NoContent(http.StatusConflict)
	}
	var p []string
	if err := c.Bind(&p); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := s.s.SetSession(sessionID, p); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

// EndSession is to end a session. Remove all relevant messages
func (s *Server) EndSession(c echo.Context) error {
	sessionID := c.Param("sessionID")
	if strings.Trim(sessionID, " ") == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	p, err := s.s.GetSession(sessionID)
	if err == storage.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	// delete all messages
	for _, participantID := range p {
		if err := s.s.DeleteMessage(sessionID, participantID); err != nil {
			c.Logger().Errorf("fail to delete messages of session %s, participant ID: %s,err: %w", sessionID, participantID, err)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	if err := s.s.DeleteSession(sessionID); err != nil { // delete session
		c.Logger().Errorf("fail to delete session %s,err: %w", sessionID, err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) GetMessage(c echo.Context) error {
	sessionID := c.Param("sessionID")
	participantID := c.Param("participantID")
	if strings.Trim(sessionID, " ") == "" || strings.Trim(participantID, " ") == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	c.Logger().Debug("session ID is ", sessionID, ", participant ID is ", participantID)
	messages, err := s.s.GetMessage(sessionID, participantID)
	if err == storage.ErrNotFound {
		return c.NoContent(http.StatusOK)
	}
	// delete the message after receive it
	if err := s.s.DeleteMessage(sessionID, participantID); err != nil {
		c.Logger().Errorf("fail to delete messages of session %s, participant ID: %s,err: %w", sessionID, participantID, err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, messages)
}
func (s *Server) PostMessage(c echo.Context) error {
	sessionID := c.Param("sessionID")
	if strings.Trim(sessionID, " ") == "" {
		c.Logger().Error("session ID is empty")
		return c.NoContent(http.StatusBadRequest)
	}
	c.Logger().Debug("session ID is ", sessionID)
	var m model.Message
	if err := c.Bind(&m); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusBadRequest)
	}
	p, err := s.s.GetSession(sessionID)
	if err == storage.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else {
		if err != nil {
			c.Logger().Error(err)
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	if m.To == nil {
		return c.NoContent(http.StatusBadRequest)
	}
	for _, participantID := range m.To {
		if !contains(p, participantID) {
			continue
		}

		if err := s.s.SetMessage(sessionID, participantID, m); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	}
	return c.NoContent(http.StatusCreated)
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
