package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"meteor-server/pkg/db"
)

var capes      []byte
var capeOwners []byte

func CapesHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", capes)
}

func CapeOwnersHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", capeOwners)
}

func UpdateCapes() {
	// Capes
	sb := strings.Builder{}

	for i, cape := range db.GetAllCapes() {
		if i > 0 {
			sb.WriteRune('\n')
		}

		sb.WriteString(cape.ID)
		sb.WriteRune(' ')
		sb.WriteString(cape.Url)
	}

	capes = []byte(sb.String())

	// Cape owners
	sb = strings.Builder{}

	i := 0
	for _, account := range db.GetAccountsWithCape() {
		if len(account.McAccounts) > 0 {
			cape := account.Cape
			if cape == "custom" {
				cape = account.ID
			}

			for _, uuid := range account.McAccounts {
				if i > 0 {
					sb.WriteRune('\n')
				}

				sb.WriteString(uuid.String())
				sb.WriteRune(' ')
				sb.WriteString(cape)

				i++
			}
		}
	}

	capeOwners = []byte(sb.String())
}
