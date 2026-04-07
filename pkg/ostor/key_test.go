package ostor_test

func (s *OstorTestSuite) TestGenerateCredentials() {
	user, resp, err := s.client.GenerateCredentials("user@example.org")
	s.Require().NotNil(user)
	s.Assert().NotNil(resp)
	s.Assert().NoError(err)

	s.Assert().Len(user.AccessKeys, 2)
}
