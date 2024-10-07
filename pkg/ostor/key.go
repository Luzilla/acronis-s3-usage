package ostor

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func (o *Ostor) GenerateCredentials(email string) (*OstorCreateUserResponse, *resty.Response, error) {
	var user *OstorCreateUserResponse
	resp, err := o.post(qUsers, qUsers+"&emailAddress="+email+"&genKey", &user)
	return user, resp, err
}

func (o *Ostor) RevokeKey(email, accessKeyID string) (*resty.Response, error) {
	return o.post(qUsers, qUsers+"&emailAddress="+email+"&revokeKey="+accessKeyID, nil)
}

// RotateKey attempts to rotate the accessKeyID for the given user (email).
//
// This feature is not a native feature in the APIs, so we build around it using
// our own methods, by checking the user account and verifying what can happen.
func (o *Ostor) RotateKey(email, accessKeyID string) (*AccessKeyPair, *resty.Response, error) {
	user, userResp, err := o.GetUser(email)
	if err != nil {
		if userResp.StatusCode() == http.StatusNotFound {
			return nil, nil, fmt.Errorf("user %q does not exist", email)
		}
		return nil, userResp, err
	}

	if len(user.AccessKeys) == 0 {
		return nil, nil, fmt.Errorf("user %q has no access keys to rotate", email)
	}

	foundKey := false

	// each account can have up to 2 key pairs
	var otherKeyPair *AccessKeyPair
	for _, kP := range user.AccessKeys {
		if kP.AccessKeyID == accessKeyID {
			foundKey = true
			continue
		}

		otherKeyPair = &kP
	}

	if !foundKey {
		return nil, nil, fmt.Errorf("user %q has no access key %q that could be rotated", email, accessKeyID)
	}

	revokeResp, err := o.RevokeKey(email, accessKeyID)
	if err != nil {
		return nil, revokeResp, err
	}

	genUser, genResp, err := o.GenerateCredentials(email)
	if err != nil {
		return nil, genResp, fmt.Errorf("failed to generate new credentials for %q: %w", email, err)
	}

	if len(genUser.AccessKeys) == 0 {
		return nil, nil, fmt.Errorf("no credentials were generated (ostor error) for %q", email)
	}

	if otherKeyPair == nil { // no other key was there, no need to filter
		return &genUser.AccessKeys[0], genResp, nil
	}

	for _, kP := range genUser.AccessKeys {
		if kP.AccessKeyID != otherKeyPair.AccessKeyID {
			return &kP, genResp, nil
		}
	}

	// this should not happen
	return nil, genResp, fmt.Errorf("unable to find new key pair")
}
