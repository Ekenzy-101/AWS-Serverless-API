package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Ekenzy-101/Serverless-Ecommerce-API/entity"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /accounts/register", func() {
	type Response struct {
		User    *entity.User `json:"user"`
		Message string       `json:"message"`
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
			"email":    email,
			"name":     name,
			"password": password,
		})
		if err != nil {
			return nil, err
		}

		url := baseUrl + "/accounts/register"
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
		id = ""
	})

	AfterEach(func() {
		if id != "" {
			Expect(repo.DeleteUser(ctx, id)).To(Succeed())
		}
	})

	It("should register user if inputs are valid", func() {
		res, err := execute()
		id = body.User.ID

		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(HaveHTTPStatus(http.StatusCreated))
		Expect(body).To(And(
			HaveField("User.Email", email),
			HaveField("User.Name", name),
			HaveField("User.Password", ""),
		))
	})

	It("should not register user if he already exists", func() {
		user := &entity.User{
			Email:    email,
			Name:     name,
			Password: password,
		}
		Expect(repo.RegisterUser(ctx, user)).To(Succeed())
		Expect(repo.CreateUser(ctx, user)).To(Succeed())

		id = user.ID
		res, err := execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(HaveHTTPStatus(http.StatusUnprocessableEntity))
		Expect(body).To(HaveField("Message", ContainSubstring("exists")))
	})

	It("should not register user if email input is invalid", func() {
		email = strings.Repeat("a", 256)
		res, err := execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(HaveHTTPStatus(http.StatusUnprocessableEntity))
		Expect(body).To(HaveField("Message", ContainSubstring("Email")))
	})

	It("should not register user if name input is invalid", func() {
		name = "1234"
		res, err := execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(HaveHTTPStatus(http.StatusUnprocessableEntity))
		Expect(body).To(HaveField("Message", ContainSubstring("Name")))
	})

	It("should not register user if password input is invalid", func() {
		password = "invalidpass"
		res, err := execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(HaveHTTPStatus(http.StatusUnprocessableEntity))
		Expect(body).To(HaveField("Message", ContainSubstring("Password")))
	})
})
