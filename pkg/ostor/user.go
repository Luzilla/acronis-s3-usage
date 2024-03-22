package ostor

import (
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"
)

// query parameter for user management
const qUsers string = "ostor-users"

func (o *Ostor) CreateUser(email string) error {
	_, err := o.putRequest(qUsers, qUsers+"&emailAddress="+email)
	return err
}

func (o *Ostor) ListUsers() (*OstorUsersListResponse, error) {
	var users *OstorUsersListResponse
	_, err := o.getRequest(qUsers, qUsers, &users)
	return users, err
}

func (o *Ostor) GetUser(email string) (*OstorUser, error) {
	var user *OstorUser
	resp, err := o.getRequest(qUsers, qUsers+"&emailAddress="+email, &user)

	b, _ := io.ReadAll(resp.RawResponse.Body)
	fmt.Printf("show: %s", b)
	return user, err
}

func (o *Ostor) GenerateCredentials(email string) (*resty.Response, error) {
	return o.postRequest(qUsers, qUsers+"&emailAddress="+email+"&genKey")
}

func (o *Ostor) RevokeKey(email, accessKeyID string) (*resty.Response, error) {
	return o.postRequest(qUsers, qUsers+"&emailAddress="+email+"&revokeKey="+accessKeyID)
}
