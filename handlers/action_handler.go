package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Url struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

type Actions struct {
	ctx context.Context
}

var client *redis.Client

var domain = "http://127.0.0.1:8001/"

func NewAction() *Actions {
	return &Actions{
		ctx: context.TODO(),
	}
}

func randomString(length int) string {
	return uuid.NewString()[:length]
}

func (a *Actions) RedisConnection() bool {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(a.ctx).Result()
	if err != nil {
		fmt.Println(err)
	}
	return true
}

func (a *Actions) GetUrlHandler(f *fiber.Ctx) error {

	url := f.Params("url")
	key := "url_id:" + url
	res, err := client.HGet(a.ctx, string(key), string("long")).Result()
	if err != nil {
		panic(err.Error())
	}

	return f.Redirect(res, http.StatusPermanentRedirect)
}

func (a *Actions) CreateUrlHandler(f *fiber.Ctx) error {

	payload := &struct {
		Url string `json:"url"`
	}{}

	if err := f.BodyParser(payload); err != nil {
		panic(err)
	}

	if _, err := url.ParseRequestURI(payload.Url); err != nil {
		return f.Render("result", fiber.Map{
			"Error": "Please enter a valide URL",
		})
	}
	random := randomString(8)
	key := "url_id:" + random

	err := client.HMSet(a.ctx, key, map[string]string{
		"short": key,
		"long":  payload.Url,
	}).Err()

	if err != nil {
		panic(err)
	}

	return f.Render("result", fiber.Map{
		"Data": domain + random,
	})
}
