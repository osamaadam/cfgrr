package gh

import (
	"os/user"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func GenHTTPAuth(username, password string) *http.BasicAuth {
	auth := &http.BasicAuth{
		Username: username,
		Password: password,
	}

	return auth
}

func GenAccessTokenAuth(username, token string) *http.BasicAuth {
	if username == "" {
		user, _ := user.Current()
		username = user.Username
	}
	auth := &http.BasicAuth{
		Username: username,
		Password: token,
	}

	return auth
}

func GenSSHAuth(privateKeyPath, password string) *ssh.PublicKeys {
	user, _ := user.Current()
	auth, _ := ssh.NewPublicKeysFromFile(user.Username, privateKeyPath, password)

	return auth
}
