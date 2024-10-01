package ostor

import "github.com/go-resty/resty/v2"

// query parameter for user management
const qUsers string = "ostor-users"

func (o *Ostor) CreateUser(email string) (*OstorCreateUserResponse, *resty.Response, error) {
	var user *OstorCreateUserResponse
	resp, err := o.put(qUsers, qUsers+"&emailAddress="+email, &user)
	return user, resp, err
}

func (o *Ostor) ListUsers() (*OstorUsersListResponse, error) {
	var users *OstorUsersListResponse
	_, err := o.get(qUsers, map[string]string{}, &users)
	return users, err
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

	return o.put(qUsers, params, nil)
}
