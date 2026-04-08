package ostor_test

func (s *OstorTestSuite) TestGetBuckets() {
	buckets, resp, err := s.client.GetBuckets(s.T().Context(), "user@example.org")
	s.Require().NotNil(buckets)
	s.Require().NotNil(resp)
	s.Require().NoError(err)

	s.Assert().Len(buckets.Buckets, 2)
}
