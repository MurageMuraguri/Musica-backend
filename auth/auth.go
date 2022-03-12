package auth

import (
	"musica/db"
	"musica/model"
	"musica/utils"

	"github.com/mileusna/useragent"
	//"log"
	//"net/http"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/lithammer/shortuuid/v4"
	"golang.org/x/crypto/bcrypt"
)



var d = db.Db()

type response struct {
	model.User
	Token string `json:"token"`
}
type CustomClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}
type CustomClaimsOTT struct {
	UserId string `json:"user_id"`
	Type   string `json:"type"`
	jwt.StandardClaims
}

var signingKey = []byte(utils.Env("KEY"))

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}




//TODO Shorten username to first 5 characters + random
func Register(c *gin.Context) {

	
	var mdl model.User
	var result model.User
	var user model.User
	var count int64
	c.Bind(&user)
	d.Where("email = ?", user.Email).First(&result)
	if result.ID == "" {
		//fmt.Println(result)
		var apiKey = uuid.New()
		user.ID = uuid.New()
		user.CreatedAt = time.Now()
		user.Password, _ = HashPassword(user.Password)
		d.Model(&mdl).Count(&count)
		user.UserName = strings.Split(user.Email, "@")[0] + fmt.Sprint(count)
		user.ApiKey = apiKey
		user.Enabled = true
		user.Provider = "email"

		d.Create(&user)
		var resp response
		resp.User = user
		resp.Password = ""
		resp.Token = GenerateAccessToken(resp.ID)
		c.JSON(200, resp)
	} else {
		c.JSON(200, gin.H{
			"message": "User exists",
		})
	}

}
func GoogleLogin(c *gin.Context) {
	type GToken struct {
		Token string
	}

	var result model.User
	var user model.User
	var tken GToken
	c.Bind(&tken)

	d.Where("email = ?", user.Email).First(&result)

	log.Default()
	if result.ID == "" {
		//fmt.Println(result)
		user.ID = uuid.New()
		user.CreatedAt = time.Now()
		user.Password, _ = HashPassword(user.Password)
		user.ApiKey = uuid.New()
		user.Enabled = true
		user.Provider = "email"
		d.Create(&user)
		c.JSON(200, gin.H{
			"token": GenerateAccessToken(user.ID),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "User exists",
		})
	}

}
func CheckToken(c *gin.Context) (jwt.MapClaims, bool, string) {
	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) < 2 {
		return nil, false, "Missing Token"
	} else {
		reqToken = splitToken[1]
		return ValidateAccessToken(reqToken)
	}

}
func ParseUserAgent(user_agent string) ua.UserAgent {
	return ua.Parse(user_agent)
}

func Login(c *gin.Context) {
	var result model.User
	var user model.User
	var login model.Login

	c.Bind(&user)
	d.Where("email = ?", user.Email).First(&result)
	if result.ID != "" {
		if CheckPasswordHash(user.Password, result.Password) {
			var resp response
			resp.User = result
			resp.Password = ""
			resp.Token = GenerateAccessToken(result.ID)

			//create login entry

			login.UserId = resp.ID
			login.ID = uuid.New()
			login.CreatedAt = time.Now()
			login.IpAddress = c.ClientIP()
			login.UserAgent = c.Request.UserAgent()

			d.Create(&login)

			c.JSON(200, resp)
		} else {
			c.JSON(200, gin.H{
				"message": "Wrong Password",
			})
		}
	} else {
		c.JSON(200, gin.H{
			"message": "User does not exist",
		})
	}

}
func GenerateAccessToken(uid string) string {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	// Create the Claims
	claims := CustomClaims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 360).Unix(),
			Issuer:    "trillwave.com",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		log.Fatalln(err)
	}
	return tokenString
}
func ValidateAccessToken(tokenString string) (jwt.MapClaims, bool, string) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	msg := ""
	if token.Valid {
		msg = "Valid Token"
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			msg = "Invalid token"
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			msg = "Expired Token"
		} else {
			msg = "Couldn't handle this token"
		}
	} else {
		msg = "Couldn't handle this token"
	}

	return token.Claims.(jwt.MapClaims), token.Valid, msg

}

// generates one time access tokens
func GenerateOTT(uid string) string {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	// Create the Claims
	claims := CustomClaimsOTT{
		uid,
		"ott",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "trillwave.com",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		log.Fatalln(err)
	}
	return tokenString
}


