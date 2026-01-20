package middleware

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dbconfig "golang.elasticsearch/dbconfig"
	utils "golang.elasticsearch/utils"

	"golang.elasticsearch/dto"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	qrcode "github.com/skip2/go-qrcode"
)

// @Summary MFA Activation
// @Description Multi-Factor Authenticator
// @Tags MultiFactor Authenticator
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User Id"
// @Param body body dto.MfaActivation true "Enable MFA"
// @Success 200 {array} dto.MfaActivation
// @Router /api/mfa/activate/{id} [patch]
func MfaActivate(c *gin.Context) {
	id := c.Param("id")
	var user dto.MfaActivation
	err := json.NewDecoder(c.Request.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	if user.TwoFactoEnabled {
		user, err := utils.GetUserid(id)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		if len(user) > 0 {
			key, err := totp.Generate(totp.GenerateOpts{
				Issuer:      "BARCLAYS BANK", // The name of your application
				AccountName: user[0].Email,   // The user's account identifier
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
				return
			}
			// The key.Secret() is the base32 encoded secret you must save
			secret := key.Secret()
			// The key.URL() is the otpauth URI, which can be converted into a QR code
			qrCodeURL := key.URL()

			pngBytes, err := qrcode.Encode(qrCodeURL, qrcode.Medium, 256)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate QR code: %v", err)})
				return
			}
			// 3. Base64 encode the PNG bytes
			var mfaData dto.MfaData
			// "data:image/png;base64,
			base64Encoded := base64.StdEncoding.EncodeToString(pngBytes)
			mfaData.Secret = secret
			mfaData.Qrcodeurl = string(base64Encoded)

			updateData := map[string]interface{}{
				"doc": map[string]interface{}{
					"secret":    mfaData.Secret,
					"qrcodeurl": mfaData.Qrcodeurl,
				},
			}
			payload, _ := json.Marshal(updateData)
			esClient := dbconfig.Connection()

			res, err := esClient.Update(
				"users", // Index name
				id,      // Document ID
				bytes.NewReader(payload),
				esClient.Update.WithContext(context.Background()),
			)

			if err != nil {
				c.JSON(500, gin.H{"message": "Error updating database"})
				return
			}
			defer res.Body.Close()

			if res.IsError() {
				c.JSON(500, gin.H{"message": fmt.Sprintf("Elasticsearch error: %s", res.String())})
				return
			}

			c.JSON(200, gin.H{
				"qrcodeurl": base64Encoded,
				"message":   "Multi-Factor Authenticator has been enabled."})

		}
		log.Print("MFA is enabled..............")
	} else {
		log.Print("MFA is disabled...........................")
		updateData := map[string]interface{}{
			"script": map[string]interface{}{
				"source": "ctx._source.secret = null; ctx._source.qrcodeurl = null;",
				"lang":   "painless",
			},
		}
		payload, _ := json.Marshal(updateData)
		esClient := dbconfig.Connection()

		res, err := esClient.Update(
			"users", // Index name
			id,      // Document ID
			bytes.NewReader(payload),
			esClient.Update.WithContext(context.Background()),
		)

		if err != nil {
			c.JSON(500, gin.H{"message": "Error updating database"})
			return
		}
		defer res.Body.Close()

		if res.IsError() {
			c.JSON(500, gin.H{"message": fmt.Sprintf("Elasticsearch error: %s", res.String())})
			return
		}

		c.JSON(200, gin.H{
			"message": "Multi-Factor Authenticator has been disabled."})

	}

}
