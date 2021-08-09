package rest_test

import (
	"testing"

	"github.com/dchest/uniuri"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const baseDBURI = "postgres://molepatrol:molepatrol@localhost:5432"

// nolint
const jwkPublicSet = `{
	"keys": [
		{
			"use": "sig",
			"kty": "RSA",
			"kid": "KWZ5ZSZZ",
			"alg": "RS256",
			"n": "uJOjx3H9Rfih6sMSQ5jZhK_Ko7-zdGP5syitCjgmsM_xvKX9Qes_4QY8IMP23djA1D8v1751KrCXz_IwXB5N4v4ny-VX9lDAKInb8DpVcmpfmB8sm5hYfDQ8X15sah_sfvtC9F1Ddnv7IHGJvhLhiJ_QPRFDh38PlDWpMZ3Py2sXs0g88BrLiCk2OyOekDHxu0XB3fSlZRVQReQ9Wnp_usudKhQlMxA7X2WrvsAG-yBTU6Y72mz5ni5qV9WcAO848vepfRq_QR6myifSe_aGHqLB_IWZdJtstktR4pwWqAq2gnyXS_jzaD5pqVzB0cmyeLmhFFHoMXTsygj8lpe7vQ",
			"e": "AQAB"
		}
	]
}`

// nolint
const jwkKeyPair = `{
	"use": "sig",
	"kty": "RSA",
	"kid": "KWZ5ZSZZ",
	"alg": "RS256",
	"n": "uJOjx3H9Rfih6sMSQ5jZhK_Ko7-zdGP5syitCjgmsM_xvKX9Qes_4QY8IMP23djA1D8v1751KrCXz_IwXB5N4v4ny-VX9lDAKInb8DpVcmpfmB8sm5hYfDQ8X15sah_sfvtC9F1Ddnv7IHGJvhLhiJ_QPRFDh38PlDWpMZ3Py2sXs0g88BrLiCk2OyOekDHxu0XB3fSlZRVQReQ9Wnp_usudKhQlMxA7X2WrvsAG-yBTU6Y72mz5ni5qV9WcAO848vepfRq_QR6myifSe_aGHqLB_IWZdJtstktR4pwWqAq2gnyXS_jzaD5pqVzB0cmyeLmhFFHoMXTsygj8lpe7vQ",
	"e": "AQAB",
	"d": "Y-CIRFlbSuyieU17aZahRZp2VatbKQUcTiUZlakSzqSHU6SiaXQqCdL84GIKCLvMhE14zw6hiisqyvxrzL0dOlJ7KGr-8St6_7SxjcmTCSmkdsWPttZ2Myd078pBch-6MnA2J9L7uXaXSlQFzBOddPe9j_3yg7Rusq1i05VoptJto781KvUfyWlHPvRjxAj0tKYrvy-NBxRdwEOGB_FfWajBycqV4BoBCRrnLGU7sSO2pX3w8OJ8ki6hPNCGCO0HFcjFp6_Du_Vd5IbKZFyJv3iiLKYa1tAffk873EqF-TlkkrJEUhyhYNydsdTvY7KQq7Wt0oyvaySz_I4rx6GbwQ",
	"p": "yy4kiA8etnquXSZBSkW58MbMT8shVVJyUs7gZ3GGr2luLCcTNjpgpc1H8RISTM2E1qXEnN4Z7NtMp9gPdUznGFPlRqzOrBxcHR584Ve_xe1Q36gy5REDUPkYRqtohWrFBzdaURifzat9lSumKLkJ-HBg8nxudB9I6Fkko_S3Y00",
	"q": "6I9ptTH96uuAJE6SdJH8wllojthWLkM86uh2qat7Qf_HBwkx1gLX9T2ZVaU3-x2fPgpQT5SHkDI5GqbbuNZT6QQxlRqKFqADSMC0oVJ6x6RzGa5YupN9uorJCra39xzdSb0n5cTsODMOCci_gSIIWi_NQvssDZKWDaDukMvDojE",
	"dp":"rINAN1oHLM8bjzG1C_gJ-ZsBzNpfMg-vzAmlVY952SQ-jDSdRlTozL5w0AoVCasSmCHlv3-BKa_F9VkpPuKN8QUCzjBZxp7Jw1uokrirtsVZ_pzUodQBKdZmO1K8i1NteUQRZnvu63UpSufly_vYsF3SovDt46DQiZ9u0dstfp0",
	"dq":"1QlC-XRxOTlAaoH7kYOGOnby3B_7WmfHrx0CTs1CnSP0q1JV78ktEX-7LgSqsoPhM1D5Xt0eDg6j1vFRSBI2TbfEv-TO6IjuWUAGd915kdbohXb72vZvb2nhXsog3eL4J6t6l_X7ukOysW3PWDjX094ENz6ljU1h3dw1jjjK3pE",
	"qi":"Yo3gbzeFja4waIcEIxXlhZ_rKFoXPpkwaPubb1eNNKZpJ9QQqWTJ4abq4eb6T3I6jxYMToZRXHYyh9b7QcrnVXmx7y6QurKUyj7OQAz0AfGEbzn8lj0WMHvUpUMIeqLSEB3DwpqyNJQO7E7EmRk6UZMRzLZC1kuk2_5GWC2RZHY"
}`

func TestRest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "molepatrol/rest")
}

func prepareDB(dbConn *pgx.ConnPool, dbName *string) func() {
	return func() {
		setupConfig, err := pgx.ParseURI(baseDBURI + "/molepatrol?sslmode=disable")
		Expect(err).To(BeNil())
		setupConn, err := pgx.Connect(setupConfig)
		Expect(err).To(BeNil())

		*dbName = uniuri.NewLenChars(16, []byte("abcdefghijklmnopqrstuvwxyz"))
		_, err = setupConn.Exec("CREATE DATABASE " + *dbName)
		Expect(err).To(BeNil())
		Expect(setupConn.Close()).To(BeNil())

		testConfig, err := pgx.ParseURI(baseDBURI + "/" + *dbName + "?sslmode=disable")
		Expect(err).To(BeNil())

		conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
			ConnConfig: testConfig,
		})
		Expect(err).To(BeNil())

		*dbConn = *conn // nolint

		m, err := migrate.New(
			"file://../../migrations",
			baseDBURI+"/"+*dbName+"?sslmode=disable",
		)
		Expect(err).To(BeNil())
		Expect(m.Up()).To(BeNil())
		srcErr, dbErr := m.Close()
		Expect(srcErr).To(BeNil())
		Expect(dbErr).To(BeNil())
	}
}

func cleanupDB(dbConn *pgx.ConnPool, dbName *string) func() {
	return func() {
		dbConn.Close()

		setupConfig, err := pgx.ParseURI(baseDBURI + "/molepatrol?sslmode=disable")
		Expect(err).To(BeNil())
		setupConn, err := pgx.Connect(setupConfig)
		Expect(err).To(BeNil())

		_, err = setupConn.Exec("DROP DATABASE " + *dbName)
		Expect(err).To(BeNil())
		Expect(setupConn.Close()).To(BeNil())
	}
}
