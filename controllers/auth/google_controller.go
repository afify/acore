package auth

import (
	authModel "acore/models/auth"
	"acore/models/session"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	Scopes:       []string{"openid", "email", "profile"},
	Endpoint:     google.Endpoint,
}

type userInfo struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := uuid.NewString()
	setOAuthStateCookie(w, state)
	url := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if err := verifyState(r); err != nil {
		slog.Error("GoogleCallback:", "verifyState", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// 2) Exchange code for OAuth2 token
	tok, err := exchangeCode(r)
	if err != nil {
		slog.Error("GoogleCallback:", "exchangeCode", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// 3) Fetch userinfo from Google
	ui, err := fetchUserInfo(r, tok)
	if err != nil {
		slog.Error("GoogleCallback:", "fetchUserInfo", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// 4) Lookup or create a local User + provider link
	userID, err := findOrCreateUser(ui)
	if err != nil {
		slog.Error("GoogleCallback:", "findOrCreateUser", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// 5) Issue our own session cookie
	if err := session.CreateSession(w, r, userID, session.SessionTypeWeb, authModel.AuthProviderGoogle); err != nil {
		slog.Error("GoogleCallback:", "CreateSession", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// 6) Finally, send them home
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func verifyState(r *http.Request) error {
	c, err := r.Cookie(authModel.OauthStateCookieName)
	if err != nil {
		slog.Error("verifyState: missing state cookie", "error", err)
		return errors.New("oauth state cookie not found")
	}

	googleState := r.URL.Query().Get("state")
	if googleState != c.Value {
		slog.Error("verifyState: state mismatch",
			"cookie", c.Value,
			"query", googleState,
		)
		return errors.New("invalid oauth state")
	}
	return nil
}

func exchangeCode(r *http.Request) (*oauth2.Token, error) {
	code := r.URL.Query().Get("code")
	return googleOAuthConfig.Exchange(r.Context(), code)
}

func fetchUserInfo(r *http.Request, tok *oauth2.Token) (*userInfo, error) {
	client := googleOAuthConfig.Client(r.Context(), tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ui userInfo
	if err := json.NewDecoder(resp.Body).Decode(&ui); err != nil {
		return nil, err
	}
	return &ui, nil
}

func findOrCreateUser(ui *userInfo) (uuid.UUID, error) {
	// 1) Try lookup by provider
	id, err := authModel.GetUserByProvider(authModel.AuthProviderGoogle, ui.Sub)
	if err != nil {
		// 1a) Not linked yet → create user + link
		if errors.Is(err, sql.ErrNoRows) {
			slog.Info("No existing provider link, creating user")

			u, err := authModel.CreateUser(authModel.SignupReq{
				UserName: ui.Email,
				Email:    ui.Email,
				Password: "", // no password for OAuth users
			})
			if err != nil {
				return uuid.Nil, err
			}

			if err := authModel.LinkProvider(u.ID, authModel.AuthProviderGoogle, ui.Sub); err != nil {
				return uuid.Nil, err
			}

			return u.ID, nil
		}
		// 1b) Some other error
		return uuid.Nil, err
	}

	// 2) Found an existing user
	return id, nil
}

func setOAuthStateCookie(w http.ResponseWriter, state string) {
	http.SetCookie(w, &http.Cookie{
		Name:     authModel.OauthStateCookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(5 * time.Minute),
		SameSite: http.SameSiteNoneMode,
	})
}
