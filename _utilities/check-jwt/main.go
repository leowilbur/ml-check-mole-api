package main

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/leowilbur/ml-check-mole-api/pkg/auth"
)

const input = `eyJraWQiOiJMdXhTcDlOZk8rdHRINkM4VDA3Sk9ZM3dlVStwdFJlNGN3blpxVmdnK3R3PSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiIxMDY5ZDcyZC0zMDIyLTQ2YmQtYTA2MC04MTY4YTc0MjNlMzkiLCJkZXZpY2Vfa2V5IjoiYXAtc291dGhlYXN0LTJfMGJlZTlhMTQtYzlhMy00OTRkLTgxNjMtNWUzNmFhZmMxMjAxIiwiY29nbml0bzpncm91cHMiOlsiRG9jdG9ycyJdLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAuYXAtc291dGhlYXN0LTIuYW1hem9uYXdzLmNvbVwvYXAtc291dGhlYXN0LTJfZ2ZTdXVIdzZlIiwiY2xpZW50X2lkIjoibHUyazdwdGZqYjRnaDA1bmtnMWtkM3U5diIsImV2ZW50X2lkIjoiZjgxZTNlOTYtZTY5NC0xMWU4LTliYTAtYzM5NGYwZDVlMWI2IiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU0MjAzODgxMiwiZXhwIjoxNTQyMDQyNDEyLCJpYXQiOjE1NDIwMzg4MTIsImp0aSI6ImUxYzQwZjgzLTNkNTYtNDU5MS1hZTVhLTJmZDIzMzBiYmZjOCIsInVzZXJuYW1lIjoiMTA2OWQ3MmQtMzAyMi00NmJkLWEwNjAtODE2OGE3NDIzZTM5In0.OWdfSAVS9csd5sMDuAL6dLbNsXm5hx1ihERvHWuMwPWGzBovNVBGAJViAlSFU7w-4M01z2d8cHnEFSQcZPhseZc1hsnTJX8tXU5WBzgmKR7BJBGDzxzKrFNoXaQ1EWQVFYwUy6VDgV3giWFNgB1vNsifBBqjQcwpVln-ohdd8IfyN2OzWAIQ4-yjSXL6UFVjIIVlKrOSjgr6Lh0xYCwOllPia3rhWmTxJyhljI5GeHt-YkO4VrqHEiqWC1GvI720fud1Lk6PoT2SsexcDntRZai5gs1yljAPlzeeuNgaTe5TpwZCkZ3yFiCJluot2Iw0kX1xRawrevKNdVYMSBP8kA`

func main() {
	keys, err := auth.CognitoJWK("")
	if err != nil {
		panic(err)
	}

	token, err := jwt.Parse(input, func(token *jwt.Token) (interface{}, error) {
		// Cognito users RS256 for JWTs
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// kid is the RSA key id
		keyID, ok := token.Header["kid"]
		if !ok {
			return nil, errors.New("Key ID missing from the JWT token")
		}

		// Cast it as string
		castedKeyID, ok := keyID.(string)
		if !ok {
			return nil, errors.New("Key ID has an invalid type")
		}

		// Match the keys in our JWK store
		matchedKeys := keys.Key(castedKeyID)
		if len(matchedKeys) == 0 {
			return nil, errors.New("Unknown key ID")
		}

		return matchedKeys[0].Key, nil
	})

	fmt.Println(err)
	fmt.Println(token)
}
