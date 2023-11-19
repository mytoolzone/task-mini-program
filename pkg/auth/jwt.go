package auth

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Auth interface {
	GenerateToken(username string, userID int) (string, error)
	ParseToken(token string) (*Claims, error)
}

// Claims Claims是一些用户信息状态和额外的jwt参数
type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	jwt.StandardClaims
}

type Jwt struct {
	jwtSecret []byte
	appName   string
}

func NewAuthJwt(jwtSecret []byte, appName string) Jwt {
	return Jwt{
		jwtSecret: jwtSecret,
		appName:   appName,
	}
}

func (j Jwt) GenerateToken(username string, userID int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Hour * 720).Unix()

	claims := Claims{
		Username: username,
		UserID:   userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime, // 过期时间
			Issuer:    j.appName,  //指定发行人
		},
	}
	// 该方法内部生成签名字符串，再用于获取完整、已签名的token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(j.jwtSecret)
	return token, err
}

func (j Jwt) ParseToken(token string) (*Claims, error) {
	// 用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.jwtSecret, nil
	})
	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目结构体都是用指针传递，节省空间
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid { // Valid()验证基于时间的声明
			return claims, nil
		}
	}
	return nil, err
}
