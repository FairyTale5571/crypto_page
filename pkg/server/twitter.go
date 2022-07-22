package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fairytale5571/crypto_page/pkg/models"
	"github.com/fairytale5571/crypto_page/pkg/storage"
	"github.com/gin-gonic/gin"
)

func (r *Router) twitterHandle(c *gin.Context) {
	verifier := c.Request.URL.Query().Get("oauth_verifier")
	tokenKey := c.Request.URL.Query().Get("oauth_token")
	id, err := r.redis.Get(fmt.Sprintf("twitter:ts_%s_id", tokenKey), storage.UserTwitter)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}

	accessToken, err := r.settings.TwitterConfig.AuthorizeTokenWithParams(models.Tokens[tokenKey], verifier, map[string]string{
		"user_id":     "",
		"screen_name": "",
	})

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}
	client, err := r.settings.TwitterConfig.MakeHttpClient(accessToken)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}
	url := "https://api.twitter.com/1.1/friends/ids.json?stringify_ids=true"
	response, err := client.Get(url)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}
	defer response.Body.Close()

	type IdsResult struct {
		IDs        []string `json:"ids"`
		NextCursor uint64   `json:"next_cursor"`
	}
	var idsResult IdsResult
	json.NewDecoder(response.Body).Decode(&idsResult)

	var ids []string
	ids = append(ids, idsResult.IDs...)
	cursor := idsResult.NextCursor
	for cursor != 0 {
		response, err = client.Get(fmt.Sprintf(url+"&cursor=%d", cursor))
		if err != nil {
			return
		}
		defer response.Body.Close()
		json.NewDecoder(response.Body).Decode(&idsResult)
		ids = append(ids, idsResult.IDs...)
		cursor = idsResult.NextCursor
	}
	const user001 = "1529780626689253377"
	for _, v := range ids {
		if v == user001 {
			fmt.Fprintln(c.Writer, "Nice, now you can go back to telegram")
			r.bot.TwitterValid(id, accessToken.AdditionalData["screen_name"])
			return
		}
	}
	r.bot.TwitterNotValid(id)
	fmt.Fprintln(c.Writer, "Try again")
}
