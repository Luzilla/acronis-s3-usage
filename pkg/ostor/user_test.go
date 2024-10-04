package ostor_test

func (s *OstorTestSuite) TestGetUser() {
	user, resp, err := s.client.GetUser("user@example.org")
	s.Assert().NotNil(user)
	s.Assert().NotNil(resp)
	s.Assert().NoError(err)
}

func (s *OstorTestSuite) TestLockUser() {
	resp, err := s.client.LockUnlockUser("user@example.org", true)
	s.Assert().NotNil(resp)
	s.Assert().NoError(err)
}

func (s *OstorTestSuite) TestListUsers() {
	users, resp, err := s.client.ListUsers(false)
	s.Assert().NotNil(users)
	s.Assert().NotNil(resp)
	s.Assert().NoError(err)
}
