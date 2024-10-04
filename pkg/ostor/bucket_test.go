package ostor_test

func (s *OstorTestSuite) TestGetBuckets() {
	buckets, resp, err := s.client.GetBuckets("user@example.org")
	s.Assert().NotNil(buckets)
	s.Assert().NotNil(resp)
	s.Assert().NoError(err)

	s.Assert().Len(buckets.Buckets, 2)
}
