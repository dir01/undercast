package server_test

func (suite *ServerSuite) TestSuccessfulLogin() {
	resp := suite.requestAPI("POST", "/api/auth/login", map[string]string{"password": suite.globalPassword})
	suite.Assert().Contains(resp.Header().Get("Set-Cookie"), "auth-session=")
	suite.Assert().Equal(
		`{"error":"","payload":"OK","status":"success"}`,
		resp.Body.String(),
	)
}

func (suite *ServerSuite) TestUnsuccessfulLogin() {
	resp := suite.requestAPI("POST", "/api/auth/login", map[string]string{"password": "wrong-password"})
	cookie := resp.Header().Get("Set-Cookie")
	suite.Assert().Empty(cookie)
	suite.Assert().Equal(
		`{"error":"wrong_password","payload":null,"status":"error"}`,
		resp.Body.String(),
	)
}

func (suite *ServerSuite) TestGetProfile() {
	suite.requestAPI("POST", "/api/auth/login", map[string]string{"password": suite.globalPassword})
	resp := suite.requestAPI("GET", "/api/auth/profile", nil)
	suite.Assert().Equal(
		`{"error":"","payload":{"isActive":true},"status":"success"}`,
		resp.Body.String(),
	)
}

func (suite *ServerSuite) TestFailGetProfile() {
	resp := suite.requestAPI("GET", "/api/auth/profile", nil)
	suite.Assert().Equal(
		`{"error":"no_profile","payload":null,"status":"error"}`,
		resp.Body.String(),
	)
}

func (suite *ServerSuite) TestLogout() {
	resp := suite.requestAPI("POST", "/api/auth/login", map[string]string{"password": suite.globalPassword})
	suite.Require().Contains(resp.Header().Get("Set-Cookie"), "auth-session=")
	resp = suite.requestAPI("POST", "/api/auth/logout", nil)
	// TODO: test that expiration date is around current date
	suite.Require().Contains(resp.Header().Get("Set-Cookie"), "Expires=")
}
