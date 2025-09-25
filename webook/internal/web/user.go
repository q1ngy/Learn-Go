package web

import (
	"net/http"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/q1ngy/Learn-Go/webook/internal/domain"
	"github.com/q1ngy/Learn-Go/webook/internal/service"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	birthdayRegexPattern = `^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`

	bizLogin = "login"
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	birthdayRexExp *regexp.Regexp
	svc            *service.UserService
	codeSvc        *service.CodeService
}

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None), // 复用编译好的正则对象
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		birthdayRexExp: regexp.MustCompile(birthdayRegexPattern, regexp.None),
		svc:            svc,
		codeSvc:        codeSvc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/users")
	group.POST("signup", h.SignUp)
	//group.POST("login", h.Login)
	group.POST("login", h.LoginJWT)
	group.POST("edit", h.Edit)
	group.GET("profile", h.Profile)

	group.POST("/login_sms/code/send", h.SendSMSLoginCode)
	group.POST("/login_sms", h.LoginSMS)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不对")
		return
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	err = h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.EmailDuplicateErr:
		ctx.String(http.StatusOK, err.Error())
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		h.setJWT(ctx, u.Id)
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) setJWT(ctx *gin.Context, uid int64) {
	uc := UserClaims{
		Uid:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	ctx.Header("x-jwt-token", tokenStr)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Options(sessions.Options{
			MaxAge: 10, // 同时影响 session 时间
		})
		sess.Set("userId", user.Id)
		err := sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或密码错误")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}
func (h *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	//sess := sessions.Default(ctx)
	//uid := sess.Get("userId").(int64)

	user := ctx.MustGet("user").(UserClaims)

	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}

	err = h.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       user.Uid,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "修改成功")

}
func (h *UserHandler) Profile(ctx *gin.Context) {
	//sess := sessions.Default(ctx)
	//uid := sess.Get("userId").(int64)

	//idVal, _ := ctx.Get("uid")
	//uid, ok := idVal.(int64)
	//if !ok {
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}

	u := ctx.MustGet("user").(UserClaims)
	uid := u.Uid

	user, err := h.svc.FindById(ctx, uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}
	ctx.JSON(http.StatusOK, User{
		Nickname: user.Nickname,
		Email:    user.Email,
		AboutMe:  user.AboutMe,
		Birthday: user.Birthday.Format(time.DateOnly),
	})
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号码",
		})
		return
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"Phone"`
		Code  string `json:"Code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码不对，请重新输入",
		})
		return
	}
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	h.setJWT(ctx, u.Id)
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})

}

var JWTKey = []byte("FBW37w0expzAWYvQMurrHpBzexyi8245")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
