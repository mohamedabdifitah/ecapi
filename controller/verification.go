package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/service"
)

func VerifyOtpPhone(c *gin.Context) {
	var body map[string]string = make(map[string]string)
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	phone, ok := body["phone"]
	if !ok {
		c.String(400, "there is no email address, please enter a valid phone number")
		return
	}
	otp, ok := body["otp"]
	if !ok {
		c.String(400, "otp is empty, please enter otp we send to you phone number")
		return
	}
	validOtp := service.RedisClient.Get(service.Ctx, "otp"+":"+phone).Val()
	if otp == validOtp {
		c.String(200, "verify otp")
		return
	}
	c.String(401, "Invalid otp , please try to enter the otp we sent you.")
}
func VerifyOtpEmail(c *gin.Context) {
	var body map[string]string = make(map[string]string)
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	email, ok := body["email"]
	if !ok {
		c.String(400, "there is no email address, please enter a valid email address")
		return
	}
	otp, ok := body["otp"]
	if !ok {
		c.String(400, "otp is empty, please enter otp we send to you email address")
		return
	}
	validOtp := service.RedisClient.Get(service.Ctx, "otp:"+email).Val()
	if otp == validOtp {
		c.String(200, "verify otp")
		return
	}
	c.String(401, "Invalid otp , please try to enter the otp we sent you.")
}
