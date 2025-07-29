package utils

// 用于放一些工具的函数
import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func HashPassword(pwd string) (string, error) {
	//cost:12表示2^12次方迭代
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
}

// GenerateJWT :有些信息会存储在浏览器上,但是害怕它在浏览器篡改之后自己再直接使用.所以jwt其实就是自己在服务器上对这个信息进行签名,然后返回给浏览器,下次浏览器再请求的时候,会把这个签名带上,
// 通过验证这个签名,来判断这个请求是否合法有无篡改
func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), //3天后过期
	})
	//生成签名的 JWT（JSON Web Token）。它的作用是将 JWT 的头部（Header）和载荷（Payload）使用指定的签名算法（如 HS256）进行签名，
	//并返回一个完整的 JWT 字符串。
	signedToken, err := token.SignedString([]byte("scretkey")) //scretkey是签名密钥
	return "Bearer" + signedToken, err
}
