package sessions

import (
	"github.com/gorilla/sessions"
)

//Store has capitalized S because it must be read outside of the package
var Store = sessions.NewCookieStore(([]byte("paZZw0rd123")))
