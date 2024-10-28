package ostor

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// query parameter for user management
const qUsers string = "ostor-users"

func (o *Ostor) CreateUser(email string) (*OstorCreateUserResponse, *resty.Response, error) {
	var user *OstorCreateUserResponse
	resp, err := o.put(qUsers, qUsers+"&emailAddress="+email, &user)
	return user, resp, err
}

func (o *Ostor) DeleteUser(email string) (*resty.Response, error) {
	resp, err := o.delete(qUsers, emailMap(email))
	if err != nil {
		return resp, err
	}

	if resp.StatusCode() != http.StatusNoContent {
		err = fmt.Errorf("wrong status code: %d", resp.StatusCode())
		return resp, err
	}

	return resp, nil
}

func (o *Ostor) ListUsers(usage bool) (*OstorUsersListResponse, *resty.Response, error) {
	var users *OstorUsersListResponse

	params := map[string]string{}
	if usage {
		params["space"] = ""
	}

	resp, err := o.get(qUsers, params, &users)
	return users, resp, err
}

func (o *Ostor) GetUser(email string) (*OstorUser, *resty.Response, error) {
	var user *OstorUser
	resp, err := o.get(qUsers, map[string]string{"emailAddress": email}, &user)
	return user, resp, err
}

func (o *Ostor) LockUnlockUser(email string, lock bool) (*resty.Response, error) {
	params := qUsers + "&emailAddress=" + email
	if lock {
		params += "&disable"
	} else {
		params += "&enable"
	}

	return o.post(qUsers, params, nil)
}
