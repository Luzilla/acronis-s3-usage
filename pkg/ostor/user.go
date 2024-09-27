package ostor

// query parameter for user management
const qUsers string = "ostor-users"

func (o *Ostor) CreateUser(email string) (*OstorCreateUserResponse, error) {
	var user *OstorCreateUserResponse
	_, err := o.put(qUsers, qUsers+"&emailAddress="+email, &user)
	return user, err
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
