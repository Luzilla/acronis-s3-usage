package ostor

import "github.com/go-resty/resty/v2"

func (o *Ostor) GenerateCredentials(email string) (*OstorCreateUserResponse, *resty.Response, error) {
	var user *OstorCreateUserResponse
	resp, err := o.post(qUsers, qUsers+"&emailAddress="+email+"&genKey", user)
	return user, resp, err
}

func (o *Ostor) RevokeKey(email, accessKeyID string) (*resty.Response, error) {
	return o.post(qUsers, qUsers+"&emailAddress="+email+"&revokeKey="+accessKeyID, nil)
}
