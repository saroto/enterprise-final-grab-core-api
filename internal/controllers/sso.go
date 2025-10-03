package controller

import (
	"enterprise_core/internal/database"
	"enterprise_core/internal/query"
	"enterprise_core/internal/utils"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)
func SyncSystem(c *gin.Context, token string) (string, []string) {
	// Create an HTTP client with a cookie jar to store and send cookies
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// Get subsystem URLs from environment
	subsystemURLs := strings.Split(os.Getenv("SUBSYSTEM_URLS"), ",")
	mainDomain := os.Getenv("MAIN_DOMAIN") // Get main domain from env

	// Set the "sso" cookie for the main domain before making requests
	c.SetCookie("sso", token, 3600, "/", mainDomain, false, true)

	var syncedURLs []string

	for _, url := range subsystemURLs {
		if url == "" {
			continue
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Println("Error creating request for URL:", url, "Error:", err)
			continue
		}

		// Set the Authorization header with the token
		req.Header.Set("Authorization", "Bearer "+token)

		// Send the HTTP request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request to URL:", url, "Error:", err)
			continue
		}


		// for _, cookie := range resp.Cookies() {
		// 	c.SetCookie(cookie.Name, cookie.Value, 3600, "/", mainDomain, false, true)
		// }

		// If the response is successful, add the URL to the syncedURLs slice
		if resp.StatusCode == http.StatusOK {
			syncedURLs = append(syncedURLs, url)
		}

		// Close the response body to free resources
		resp.Body.Close()
	}

	// Return the appropriate message and the list of synced URLs
	if len(syncedURLs) > 0 {
		return "Successfully synced to subsystems", syncedURLs
	} else {
		return "No subsystems were synced", nil
	}
}

func GetSSO(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user ID from context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		data, token, expiry, err := query.GetSSOQuery(c, db, userID.(int))
		
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		utils.ResponseWithData(c, "Get SSO Successfully", gin.H{"user": data,
			"tokens": gin.H{
				"access_token": token,
				"expiry_date": expiry,
			},
			})
	}
}