package ostor

import (
	"context"
	"fmt"
	"net/http"
)

// query parameter for user management
const qUsers string = "ostor-users"

func (o *Ostor) CreateUser(ctx context.Context, email string) (*OstorCreateUserResponse, *http.Response, error) {
	var user *OstorCreateUserResponse
	resp, err := o.put(ctx, qUsers, map[string]string{"emailAddress": email}, &user)
	return user, resp, err
}

func (o *Ostor) DeleteUser(ctx context.Context, email string) (*http.Response, error) {
	resp, err := o.delete(ctx, qUsers, emailMap(email))
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("wrong status code: %d", resp.StatusCode)
		return resp, err
	}

	return resp, nil
}

func (o *Ostor) ListUsers(ctx context.Context, usage bool) (*OstorUsersListResponse, *http.Response, error) {
	var users *OstorUsersListResponse

	params := map[string]string{}
	if usage {
		params["space"] = ""
	}

	resp, err := o.get(ctx, qUsers, params, &users)
	return users, resp, err
}

func (o *Ostor) GetUser(ctx context.Context, email string) (*OstorUser, *http.Response, error) {
	var user *OstorUser
	resp, err := o.get(ctx, qUsers, emailMap(email), &user)
	return user, resp, err
}

func (o *Ostor) LockUnlockUser(ctx context.Context, email string, lock bool) (*http.Response, error) {
	params := map[string]string{"emailAddress": email}
	if lock {
		params["disable"] = ""
	} else {
		params["enable"] = ""
	}

	return o.post(ctx, qUsers, params, nil)
}
