package main

import (
	"log"
	"net/http"

	fcm "github.com/appleboy/go-fcm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jackc/pgx"
	"github.com/namsral/flag"
	"github.com/pzduniak/gateway"

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/auth"
	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/rest"
)

var (
	bind  = flag.String("bind", ":8080", "address which we should bind to")
	debug = flag.Bool("debug", true, "start in debug mode (non-lambda)")

	postgres = flag.String("postgres", "", "dsn of the postgres database")

	authURL = flag.String(
		"auth_url", "https://auth.checkmoleapp.demo-redisys.com",
		"base URL of the auth app",
	)
	awsRegion = flag.String("aws_region", "ap-southeast-2", "aws region to use")

	fcmKey = flag.String(
		"fcm_key",
		"AAAApWHQKUU:APA91bHVh8MKPu5p_qfqIipaVaT1D2vuAq5uWN5GJI1O4za9LHbzRnFhYb"+
			"1f6Swcy-m0VyFU-8rR-OM62M6VKoAhOI6NlFCl7pcYlyfW4Wyf44GxeY4a7Smm5F2P"+
			"zroGxmi7aSRdDHv-",
		"fcm key to use",
	)
)

func main() {
	flag.Parse()

	log.Println("Connecting to PostgreSQL")
	connConfig, err := pgx.ParseDSN(*postgres)
	if err != nil {
		panic(err)
	}
	dbPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: connConfig,
	})
	if err != nil {
		panic(err)
	}

	jwkSet, err := auth.CognitoJWK("")
	if err != nil {
		panic(err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: awsRegion,
	})
	if err != nil {
		panic(err)
	}

	fcc, err := fcm.NewClient(*fcmKey)
	if err != nil {
		panic(err)
	}

	router, err := rest.New(
		dbPool,
		jwkSet,
		&rest.Config{
			AuthURL:    *authURL,
			AWSSession: sess,
			FCMClient:  fcc,
		},
	)
	if err != nil {
		panic(err)
	}

	log.Println("Starting the API")

	if !*debug {
		panic(gateway.ListenAndServe(*bind, router))
	}

	panic(http.ListenAndServe(*bind, router))
}
