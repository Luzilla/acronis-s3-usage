package ostor

import "github.com/go-resty/resty/v2"

// query parameter for user management
const qUsers string = "ostor-users"

func (o *Ostor) ListUsers() (*OstorUsersListResponse, error) {
	var users *OstorUsersListResponse
	_, err := o.getRequest(qUsers, qUsers, &users)
	return users, err
}

func (o *Ostor) GetUser(email string) (*OstorUser, error) {
	var user *OstorUser
	_, err := o.getRequest(qUsers, qUsers+"&emailAddress="+email, &user)
	return user, err
}

func (o *Ostor) GenerateCredentials(email string) (*resty.Response, error) {
	return o.postRequest(qUsers, qUsers+"&emailAddress="+email+"&genKey")
}

func (o *Ostor) RevokeKey(email, accessKeyID string) (*resty.Response, error) {
	return o.postRequest(qUsers, qUsers+"&emailAddress="+email+"&revokeKey="+accessKeyID)
}
