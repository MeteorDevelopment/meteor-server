package api

import (
	"bytes"
	"encoding/json"
	"github.com/meyskens/go-turnstile"
	"image"
	_ "image/png"
	"io/ioutil"
	"meteor-server/pkg/discord"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"meteor-server/pkg/auth"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"

	"github.com/segmentio/ksuid"
)

type capeInfo struct {
	db.Cape
	Title   string `json:"title"`
	Current bool   `json:"current"`
}

type accountInfo struct {
	db.Account
	DiscordName   string     `json:"discord_name"`
	DiscordAvatar string     `json:"discord_avatar"`
	Capes         []capeInfo `json:"capes"`
}

type mcUser struct {
	Id string
}

type passwordChangeInfo struct {
	accountId ksuid.KSUID
	time      time.Time
	data      string
}

type accountTimeInfo struct {
	accountId ksuid.KSUID
	time      time.Time
}

var changeEmailTokens = make(map[ksuid.KSUID]passwordChangeInfo)
var discordLinkTokens = make(map[ksuid.KSUID]accountTimeInfo)
var forgotPasswordTokens = make(map[ksuid.KSUID]accountTimeInfo)

var ts *turnstile.Turnstile

func SetupTurnstile() {
	ts = turnstile.New(core.GetPrivateConfig().CloudflareSecretKey)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := ts.Verify(r.FormValue("cf-token"), core.IP(r))
	if err != nil || !resp.Success {
		core.JsonError(w, "Failed to verify captcha, please try again.")
		return
	}

	err = auth.Register(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		core.JsonError(w, err.Error())
		return
	}

	core.Json(w, core.J{})
}

func ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	token, err := ksuid.Parse(r.URL.Query().Get("token"))
	if err != nil {
		core.JsonError(w, "Invalid token.")
		return
	}

	if !auth.ConfirmEmail(token) {
		core.JsonError(w, "Failed to confirm email address.")
		return
	}

	db.AccountCount++

	core.Json(w, core.J{})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	token, err := auth.Login(q.Get("name"), q.Get("password"))
	if err != nil {
		core.JsonError(w, "Wrong name or password.")
		return
	}

	core.Json(w, core.J{"token": token})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.Logout(r.Header.Get("Authorization"), core.GetAccountID(r))
	if err != nil {
		core.JsonError(w, err)
		return
	}

	core.Json(w, core.J{})
}

func newCapeInfo(cape db.Cape, account db.Account, title string) capeInfo {
	return capeInfo{
		cape,
		title,
		cape.ID == account.Cape,
	}
}

func getAccountInfo(account db.Account) accountInfo {
	info := accountInfo{account, "", "", make([]capeInfo, 0)}

	// Discord info
	if account.DiscordID != "" {
		user := discord.GetUser(account.DiscordID)

		info.DiscordName = user.Username + "#" + user.Discriminator
		info.DiscordAvatar = "https://cdn.discordapp.com/avatars/" + account.DiscordID + "/" + user.Avatar + ".jpg"
	}

	// Capes
	info.Capes = append(info.Capes, newCapeInfo(db.Cape{"", ""}, account, "None"))
	if account.Donator {
		cape, _ := db.GetCape("donator")
		info.Capes = append(info.Capes, newCapeInfo(cape, account, "Donator"))
	}
	if account.Admin {
		cape, _ := db.GetCape("moderator")
		info.Capes = append(info.Capes, newCapeInfo(cape, account, "Moderator"))
	}
	if account.CanHaveCustomCape {
		cape, err := db.GetCape("account_" + account.ID.String())
		if err == nil {
			info.Capes = append(info.Capes, newCapeInfo(cape, account, "Custom"))
		}
	}

	return info
}

func AccountInfoHandler(w http.ResponseWriter, r *http.Request) {
	account, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, "Could not get account.")
		return
	}

	core.Json(w, getAccountInfo(account))
}

func GetAccountByMcUuid(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("uuid"))
	if err != nil {
		core.JsonError(w, "Invalid UUID.")
		return
	}

	account, err := db.GetAccountUuid(id)
	if err != nil {
		core.JsonError(w, "No account linked to this UUID.")
		return
	}

	core.Json(w, getAccountInfo(account))
}

func GenerateDiscordLinkTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Get account
	account, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, "Could not get account.")
		return
	}

	token := ksuid.New()
	discordLinkTokens[token] = accountTimeInfo{account.ID, time.Now()}

	core.Json(w, core.J{"token": token})
}

func LinkDiscordHandler(w http.ResponseWriter, r *http.Request) {
	// Validate token
	tokenStr := r.URL.Query().Get("token")

	token, err := ksuid.Parse(tokenStr)
	if err != nil {
		core.JsonError(w, "Invalid token.")
		return
	}

	info, exists := discordLinkTokens[token]
	if !exists {
		core.JsonError(w, "Invalid token.")
		return
	}

	delete(discordLinkTokens, token)

	if time.Now().Sub(info.time).Minutes() > 5 {
		core.JsonError(w, "Invalid token.")
		return
	}

	account, err := db.GetAccountId(info.accountId)
	if err != nil {
		core.JsonError(w, "Invalid token.")
		return
	}

	// Link
	id := r.URL.Query().Get("id")

	err = account.LinkDiscord(id)
	if err != nil {
		core.JsonError(w, err.Error())
		return
	}

	core.Json(w, core.J{})
}

func UnlinkDiscordHandler(w http.ResponseWriter, r *http.Request) {
	// Get account
	account, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, "Could not get account.")
		return
	}

	account.UnlinkDiscord()

	core.Json(w, core.J{})
}

func McAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Get account
	account, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, "Could not get account.")
		return
	}

	if r.Method == "POST" {
		// Get Minecraft username
		username := r.URL.Query().Get("username")
		if username == "" {
			core.JsonError(w, "Invalid username 1.")
			return
		}

		// Get uuid
		req, _ := http.NewRequest("GET", "https://api.mojang.com/users/profiles/minecraft/"+username, bytes.NewReader([]byte{}))
		req.Header.Set("User-Agent", "Meteor Server")

		client := http.Client{}
		res, err := client.Do(req)
		if err != nil {
			core.JsonError(w, "Invalid username 2.")
			return
		}

		body, _ := ioutil.ReadAll(res.Body)
		var user mcUser
		_ = json.Unmarshal(body, &user)

		_ = res.Body.Close()

		id, err := uuid.Parse(user.Id)
		if err != nil {
			core.JsonError(w, "Invalid username 3. "+err.Error())
			return
		}

		// Add Minecraft account
		err = account.AddMcAccount(id)
		if err != nil {
			core.JsonError(w, err.Error())
			return
		}

		UpdateCapes()
	} else {
		// Get Minecraft UUID
		id, err := uuid.Parse(r.URL.Query().Get("uuid"))
		if err != nil {
			core.JsonError(w, "Invalid UUID.")
			return
		}

		// Remove Minecraft account
		account.RemoveMcAccount(id)
		UpdateCapes()
	}

	core.Json(w, core.J{})
}

func SelectCapeHandler(w http.ResponseWriter, r *http.Request) {
	// Get account
	account, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, "Could not get account.")
		return
	}

	id := r.URL.Query().Get("cape")

	if id == "" || (id == "donator" && account.Donator) || (id == "moderator" && account.Admin) || (strings.HasPrefix(id, "account_") && account.CanHaveCustomCape) {
		account.SetCape(id)
		UpdateCapes()

		core.Json(w, core.J{})
	} else {
		core.JsonError(w, "Cannot select this cape.")
	}
}

func UploadCapeHandler(w http.ResponseWriter, r *http.Request) {
	// Get account
	account, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, "Could not get account.")
		return
	}

	// Validate file
	formFile, header, err := r.FormFile("file")
	if err != nil {
		core.JsonError(w, "Invalid file.")
		return
	}

	if !strings.HasSuffix(header.Filename, ".png") {
		core.JsonError(w, "File needs to be a PNG.")
		return
	}

	config, _, err := image.DecodeConfig(formFile)
	if err != nil {
		core.JsonError(w, "Failed to read image file.")
		return
	}

	if config.Height*2 != config.Width {
		core.JsonError(w, "Wrong size. Width of the image must be 2 times larger than the height.")
		return
	}

	if config.Width > 1024 || config.Height > 1024 {
		core.JsonError(w, "Wrong size. Maximum image size is 1024 by 512.")
		return
	}

	// Save file
	file, err := os.Create("data/capes/account_" + account.ID.String() + ".png")
	if err != nil {
		core.JsonError(w, "Server error. Failed to create cape file. Please contact developers.")
		return
	}

	if !core.DownloadFile(formFile, file, w) {
		return
	}

	cape := db.Cape{ID: "account_" + account.ID.String(), Url: "https://meteorclient.com/" + strings.TrimPrefix(file.Name(), "data/")}
	db.InsertCape(cape)

	UpdateCapes()
	core.Json(w, core.J{})
}

func ChangeUsernameHandler(w http.ResponseWriter, r *http.Request) {
	// Get account
	account, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, "Could not get account.")
		return
	}

	// Validate username
	username := r.URL.Query().Get("username")
	if username == "" || strings.ContainsRune(username, ' ') {
		core.JsonError(w, "Invalid username.")
		return
	}

	_, err = db.GetAccountWithUsername(username)
	if err == nil {
		core.JsonError(w, "Account with this username already exists.")
		return
	}

	// Change username
	account.SetUsername(username)
	core.Json(w, core.J{})
}

func ChangeEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Get account
	account, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, "Could not get account.")
		return
	}

	// Validate email
	email := r.URL.Query().Get("email")

	if !core.IsEmailValid(email) {
		core.JsonError(w, "Invalid email.")
		return
	}

	_, err = db.GetAccountWithEmail(email)
	if err == nil {
		core.JsonError(w, "Email already in use.")
		return
	}

	// Send email
	token := ksuid.New()
	changeEmailTokens[token] = passwordChangeInfo{account.ID, time.Now(), email}

	core.SendEmail(account.Email, "Confirm new email", "To change the email to "+email+" go to https://meteorclient.com/confirmChangeEmail?token="+token.String()+" . The link is valid for 15 minutes.")
	core.Json(w, core.J{})
}

func ConfirmChangeEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Validate token
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		core.JsonError(w, "Invalid token.")
		return
	}

	token, err := ksuid.Parse(tokenStr)
	if err != nil {
		core.JsonError(w, "Invalid token.")
		return
	}

	info, exists := changeEmailTokens[token]
	if !exists {
		core.JsonError(w, "Invalid token.")
		return
	}

	delete(changeEmailTokens, token)
	if time.Now().Sub(info.time).Minutes() > 15 {
		core.JsonError(w, "Outdated token.")
		return
	}

	// Change email
	account, err := db.GetAccountId(info.accountId)
	if err != nil {
		core.JsonError(w, "Invalid token.")
		return
	}

	account.SetEmail(info.data)
	http.Redirect(w, r, "https://meteorclient.com/account", http.StatusPermanentRedirect)
}

func ConfirmChangeEmailHandlerApi(w http.ResponseWriter, r *http.Request) {
	// Validate token
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		core.JsonError(w, "Invalid token.")
		return
	}

	token, err := ksuid.Parse(tokenStr)
	if err != nil {
		core.JsonError(w, "Invalid token.")
		return
	}

	info, exists := changeEmailTokens[token]
	if !exists {
		core.JsonError(w, "Invalid token.")
		return
	}

	delete(changeEmailTokens, token)
	if time.Now().Sub(info.time).Minutes() > 15 {
		core.JsonError(w, "Outdated token.")
		return
	}

	// Change email
	account, err := db.GetAccountId(info.accountId)
	if err != nil {
		core.JsonError(w, "Invalid token.")
		return
	}

	account.SetEmail(info.data)

	// Invalidate tokens
	auth.Invalidate(account.ID)

	core.Json(w, core.J{})
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Get account
	account, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, "Could not get account.")
		return
	}

	// Validate old password
	oldPass := r.URL.Query().Get("old")
	if !account.PasswordMatches(oldPass) {
		core.JsonError(w, "Wrong old password.")
		return
	}

	// Validate new password
	newPass := r.URL.Query().Get("new")
	if newPass == "" || strings.ContainsRune(newPass, ' ') {
		core.JsonError(w, "Invalid new password.")
		return
	}

	// Change password
	err = account.SetPassword(newPass)
	if err != nil {
		core.JsonError(w, "Invalid password.")
		return
	}

	// Invalidate tokens
	auth.Invalidate(account.ID)

	core.Json(w, core.J{})
}

func ChangePasswordTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Validate token
	tokenRaw := r.URL.Query().Get("token")
	token, err := ksuid.Parse(tokenRaw)
	if err != nil {
		core.JsonError(w, "Invalid token.")
		return
	}

	info, exists := forgotPasswordTokens[token]
	if !exists {
		core.JsonError(w, "Invalid token.")
		return
	}

	delete(forgotPasswordTokens, token)
	if time.Now().Sub(info.time).Minutes() > 15 {

		core.JsonError(w, "Outdated token.")
		return
	}

	// Get account
	account, err := db.GetAccountId(info.accountId)
	if err != nil {
		core.JsonError(w, "Invalid token.")
		return
	}

	// Validate new password
	newPass := r.URL.Query().Get("new")
	if newPass == "" || strings.ContainsRune(newPass, ' ') {
		core.JsonError(w, "Invalid new password.")
		return
	}

	// Change password
	err = account.SetPassword(newPass)
	if err != nil {
		core.JsonError(w, "Invalid password.")
		return
	}

	// Invalidate tokens
	auth.Invalidate(account.ID)

	core.Json(w, core.J{})
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		core.JsonError(w, "Invalid email.")
		return
	}

	account, err := db.GetAccountWithEmail(email)
	if err != nil {
		core.JsonError(w, "Invalid email.")
		return
	}

	token := ksuid.New()
	forgotPasswordTokens[token] = accountTimeInfo{account.ID, time.Now()}

	core.SendEmail(email, "Forgot password", "To change the password to "+email+" go to https://meteorclient.com/changePassword?token="+token.String()+" . The link is valid for 15 minutes.")
	core.Json(w, core.J{})
}
