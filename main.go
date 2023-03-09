package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	gsm "github.com/tokopedia/secret/v2/provider/google-secret-manager"
)

func GetConfigFromGSMADC() map[string]interface{} {
	var result map[string]interface{}

	ctx := context.Background()
	//you need to change the projectID (secrets-management-development) based on environment
	//staging = secrets-management-staging
	//production = secrets-management-278603
	secretClient, err := gsm.NewClientFromCredentialPath(ctx, "gtl-itsec-development.json")
	if err != nil {
		log.Panic(err)
	}
	//get secret value
	var a struct{}
	paths := make(map[string]struct{})
	paths["hms/versions/latest"] = a

	resf, err := secretClient.MultiFetchKV(ctx, paths)
	if err != nil {
		log.Println(err)
		return result
	}

	json.Unmarshal(resf["hms/versions/latest"], &result)
	return result
}

func main() {
	godotenv.Load()

	router := http.NewServeMux()
	router.Handle("/api/force-login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("https://%s/oauth/token", os.Getenv("AUTH0_DOMAIN"))

		paySecret := fmt.Sprintf("{\"client_id\":\"%s\",\"client_secret\":\"%s\",\"audience\":\"%s\",\"grant_type\":\"%s\"}", os.Getenv("AUTH_CLIENT_ID"), os.Getenv("AUTH_CLIENT_SECRET"), os.Getenv("AUTH0_AUDIENCE"), "client_credentials")
		payload := strings.NewReader(paySecret)

		req, _ := http.NewRequest("POST", url, payload)
		req.Header.Add("content-type", "application/json")
		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(res)
		fmt.Println(string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
	router.Handle("/api/private", EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS Headers.
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Hello from a private endpoint! You need to be authenticated to see this."}`))
		}),
	))

	log.Print("Server listening on http://localhost:3010")
	if err := http.ListenAndServe("0.0.0.0:3010", router); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
