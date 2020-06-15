package controller

func (s *Server) initializeRoutes() {

	v1 := s.Router.Group("/api/v1")
	{
		v1.GET("/ping", s.Ping)
		v1.POST("/exec", s.Exec)

	}
}
