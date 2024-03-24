package ostor

import (
	"github.com/go-resty/resty/v2"
)

// query parameter for user management
const qUsers string = "ostor-users"

func (o *Ostor) CreateUser(email string) error {
	_, err := o.put(qUsers, qUsers+"&emailAddress="+email)
	return err
}

func (o *Ostor) ListUsers() (*OstorUsersListResponse, error) {
	var users *OstorUsersListResponse
	_, err := o.get(qUsers, map[string]string{}, &users)
	return users, err
}

func (o *Ostor) GetUser(email string) (*OstorUser, error) {
	var user *OstorUser
	_, err := o.get(qUsers, map[string]string{"emailAddress": email}, &user)
	return user, err
}

func (o *Ostor) GenerateCredentials(email string) (*resty.Response, error) {
	return o.post(qUsers, qUsers+"&emailAddress="+email+"&genKey")
}

func (o *Ostor) RevokeKey(email, accessKeyID string) (*resty.Response, error) {
	return o.post(qUsers, qUsers+"&emailAddress="+email+"&revokeKey="+accessKeyID)
}
