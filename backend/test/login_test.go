package test

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /accounts/login", func() {
	type Response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IDToken      string `json:"id_token"`
		Message      string `json:"message"`
	}

	var (
		body     Response
		email    string
		name     string
		password string
		id       string
	)

	var execute = func() (*http.Response, error) {
		payload, err := json.Marshal(entity.M{
			"username": email,
			"password": password,
		})
		if err != nil {
			return nil, err
		}

		url := baseUrl + "/accounts/login"
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}

		req.Header.Add("Content-type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		return res, json.NewDecoder(res.Body).Decode(&body)
	}

	BeforeEach(func() {
		body = Response{}
		email = "test@test.com"
		name = "Test Test"
		password = "Test@123"
		user := &entity.User{
			Email:    email,
			Name:     name,
			Password: password,
		}
		Expect(repo.RegisterUser(ctx, user)).To(Succeed())
		Expect(repo.CreateUser(ctx, user)).To(Succeed())
		id = user.ID
	})

	AfterEach(func() {
		if id != "" {
			Expect(repo.DeleteUser(ctx, id)).To(Succeed())
		}
	})

	It("should login user if inputs are valid", func() {
		res, err := execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(HaveHTTPStatus(http.StatusOK))
		Expect(body).To(And(
			HaveField("AccessToken", Not(BeEmpty())),
			HaveField("IDToken", Not(BeEmpty())),
			HaveField("RefreshToken", Not(BeEmpty())),
		))
	})

	It("should not login user if he doesn't exist", func() {
		Expect(repo.DeleteUser(ctx, id)).To(Succeed())

		id = ""
		res, err := execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(HaveHTTPStatus(http.StatusUnprocessableEntity))
		Expect(body).To(HaveField("Message", ContainSubstring("exist")))
	})

	It("should not login user if email input is invalid", func() {
		email = ""
		res, err := execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(HaveHTTPStatus(http.StatusUnprocessableEntity))
		Expect(body).To(HaveField("Message", ContainSubstring("Username")))
	})

	It("should not login user if password input is invalid", func() {
		password = "invalidpass"
		res, err := execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(HaveHTTPStatus(http.StatusUnprocessableEntity))
		Expect(body).To(HaveField("Message", ContainSubstring("Password")))
	})
})
