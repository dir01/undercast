package server_test

func (s *ServerSuite) TestSuccessfulLogin() {
	resp := s.requestAPI("POST", "/api/auth/login", map[string]string{"password": s.globalPassword})
	s.Assert().Contains(resp.Header().Get("Set-Cookie"), "auth-session=")
	s.Assert().Equal(
		`{"error":"","payload":"OK","status":"success"}`,
		resp.Body.String(),
	)
}

func (s *ServerSuite) TestUnsuccessfulLogin() {
	resp := s.requestAPI("POST", "/api/auth/login", map[string]string{"password": "wrong-password"})
	cookie := resp.Header().Get("Set-Cookie")
	s.Assert().Empty(cookie)
	s.Assert().Equal(
		`{"error":"wrong_password","payload":null,"status":"error"}`,
		resp.Body.String(),
	)
}

func (s *ServerSuite) TestGetProfile() {
	s.requestAPI("POST", "/api/auth/login", map[string]string{"password": s.globalPassword})
	resp := s.requestAPI("GET", "/api/auth/profile", nil)
	s.Assert().Equal(
		`{"error":"","payload":{"isActive":true},"status":"success"}`,
		resp.Body.String(),
	)
}

func (s *ServerSuite) TestFailGetProfile() {
	resp := s.requestAPI("GET", "/api/auth/profile", nil)
	s.Assert().Equal(
		`{"error":"no_profile","payload":null,"status":"error"}`,
		resp.Body.String(),
	)
}

func (s *ServerSuite) TestLogout() {
	resp := s.requestAPI("POST", "/api/auth/login", map[string]string{"password": s.globalPassword})
	s.Require().Contains(resp.Header().Get("Set-Cookie"), "auth-session=")
	resp = s.requestAPI("POST", "/api/auth/logout", nil)
	// TODO: test that expiration date is around current date
	s.Require().Contains(resp.Header().Get("Set-Cookie"), "Expires=")
}
