package invitation

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/jpurdie/authapi"
	email "github.com/jpurdie/authapi/pkg/utl/mail"
	"github.com/labstack/echo/v4"
)

func GenerateInviteTokenHash(str string) string {
	myHash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(myHash[:])
}

func GenerateInviteToken() string {
	n := 64
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (i Invitation) Create(c echo.Context, invite authapi.Invitation) error {
	op := "Create"

	numDaysExp, err := strconv.Atoi(os.Getenv("INVITATION_EXPIRATION_DAYS"))
	if err != nil {
		numDaysExp = 7 //defaulting to 7 days
	}
	now := time.Now()
	expiresAt := now.AddDate(0, 0, numDaysExp)
	tokenStr := GenerateInviteToken()
	tokenHash := GenerateInviteTokenHash(tokenStr)
	invite.ExpiresAt = &expiresAt
	invite.TokenHash = tokenHash
	invite.TokenStr = tokenStr

	/*
		Checking if the invitation email already exists for that org
	*/
	user, err := i.udb.FetchByEmail(*i.db, invite.Email)

	if user.ID > 0 { //email already exists. Need to check if it's on the same org
		for _, tempProfile := range user.Profile {
			if tempProfile.Organization.ID == int(invite.OrganizationID) {
				return &authapi.Error{
					Op:   op,
					Code: authapi.ECONFLICT,
				} // the user exists for that organization
			}
		}
	}

	existingInvite, _ := i.idb.ViewByEmail(*i.db, invite.Email, invite.OrganizationID)
	if existingInvite.ID != 0 {
		i.idb.Delete(*i.db, invite.Email, invite.OrganizationID)
	}

	//user doesn't exist
	err = i.idb.Create(*i.db, invite)
	if err != nil {
		return err
	}

	err = email.SendInvitationEmail(&invite)
	if err != nil {
		return err
	}

	return nil
}

func (i Invitation) View(c echo.Context, tokenPlainText string) (authapi.Invitation, error) {
	tokenHash := GenerateInviteTokenHash(tokenPlainText)
	return i.idb.View(*i.db, tokenHash)
}
func (i Invitation) Delete(c echo.Context, email string, orgID int) error {
	return i.idb.Delete(*i.db, email, orgID)
}
func (i Invitation) List(c echo.Context, orgID int, includeExpired bool, includeUsed bool) ([]authapi.Invitation, error) {
	return i.idb.List(*i.db, orgID, includeExpired, includeUsed)
}

func (i Invitation) CreateUser(c echo.Context, prof authapi.Profile, invite authapi.Invitation) error {
	return i.idb.CreateUser(*i.db, prof, invite)
}
