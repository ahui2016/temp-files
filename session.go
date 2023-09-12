package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const (
	cookieName     = "TempFilesCookie"
	passwordMaxTry = 100
)

var passwordTry = 0

var SessionStore = session.New()

func sessionCheck(c *fiber.Ctx, key string) bool {
	sess, err := SessionStore.Get(c)
	val := sess.Get(key)
	if err != nil || val == nil {
		return false
	}
	return val.(bool)
}

func sessionSet(c *fiber.Ctx, key string, val any) error {
	sess, err := SessionStore.Get(c)
	if err != nil {
		return err
	}
	sess.Set(key, val)
	return sess.Save()
}

func sessionDelete(c *fiber.Ctx, key string) error {
	sess, err := SessionStore.Get(c)
	if err != nil {
		return err
	}
	sess.Delete(key)
	return sess.Save()
}

func checkLoginMidWare(c *fiber.Ctx) error {
	if isLoggedOut(c) {
		return fmt.Errorf("Require Login")
	}
	return c.Next()
}

func isLoggedIn(c *fiber.Ctx) bool {
	return sessionCheck(c, cookieName)
}

func isLoggedOut(c *fiber.Ctx) bool {
	return !isLoggedIn(c)
}

func checkPasswordTry(c *fiber.Ctx) error {
	if passwordTry >= passwordMaxTry {
		return fmt.Errorf("No more try. Input wrong password too many times.")
	}
	return nil
}

func checkPassword(c *fiber.Ctx) error {
	if passwordTry > passwordMaxTry {
		return fmt.Errorf("No more try. Input wrong password too many times.")
	}
	type Pass struct {
		Word string `json:"pwd" form:"pwd" validate:"required"`
	}
	pass := new(Pass)
	if err := c.BodyParser(pass); err != nil {
		return err
	}
	if pass.Word != app_config.Password {
		passwordTry++
		return fmt.Errorf("wrong password")
	}
	passwordTry = 0
	return nil
}
