package dengovie

import (
	"bytes"
	"database/sql"
	"dengovie/internal/domain"
	storeTypes "dengovie/internal/store/types"
	"dengovie/internal/web"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestController_RequestCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		mocks     func(e *env) (*gin.Context, *httptest.ResponseRecorder)
		respCheck func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder)
	}{
		// Существующие тесты...
		{
			name: "OK",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"telegramAlias": "alias"}`))
				ctx.Request = req

				e.storage.EXPECT().GetUserByAlias(mock.Anything, storeTypes.GetUserByAliasInput{
					Alias: "alias",
				}).Return(storeTypes.User{
					ID:     1,
					Name:   "Timr",
					Alias:  "alias",
					ChatID: 123,
				}, nil)

				e.sender.EXPECT().SendMessageToUserByAlias(
					mock.Anything,
					"alias",
					"Привет! Кое-кто запросил код для входа: `111`").
					Return(nil)

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, "", w.Body.String())
			},
		},
		{
			name: "FAIL, error from db",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"telegramAlias": "alias"}`))
				ctx.Request = req

				e.storage.EXPECT().GetUserByAlias(mock.Anything, storeTypes.GetUserByAliasInput{
					Alias: "alias",
				}).Return(storeTypes.User{}, assert.AnError)

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
				require.Equal(t, "", w.Body.String())
			},
		},
		{
			name: "OK, no user",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"telegramAlias": "alias"}`))
				ctx.Request = req

				e.storage.EXPECT().GetUserByAlias(mock.Anything, storeTypes.GetUserByAliasInput{
					Alias: "alias",
				}).Return(storeTypes.User{}, sql.ErrNoRows)

				e.storage.EXPECT().CreateUser(ctx, storeTypes.CreateUserInput{
					Name:  fmt.Sprintf("Имя %v", "alias"),
					Alias: "alias",
				}).Return(storeTypes.User{
					ID:     1,
					Name:   "Timr",
					Alias:  "alias",
					ChatID: 123,
				}, nil)

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, w.Code)
				require.Equal(t, `{"message":"А теперь напиши в бота, чтобы получить код доступа"}`, w.Body.String())
			},
		},
		{
			name: "FAIL, invalid JSON body",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"telegramAlias": "alias"`)) // Неполный JSON
				ctx.Request = req

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
				require.Equal(t, "", w.Body.String())
			},
		},
		{
			name: "FAIL, create user error after not found",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"telegramAlias": "alias"}`))
				ctx.Request = req

				e.storage.EXPECT().GetUserByAlias(mock.Anything, storeTypes.GetUserByAliasInput{
					Alias: "alias",
				}).Return(storeTypes.User{}, sql.ErrNoRows)

				e.storage.EXPECT().CreateUser(ctx, storeTypes.CreateUserInput{
					Name:  fmt.Sprintf("Имя %v", "alias"),
					Alias: "alias",
				}).Return(storeTypes.User{}, assert.AnError) // Ошибка при создании

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
				require.Equal(t, "", w.Body.String())
			},
		},
		{
			name: "FAIL, send message error for existing user",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"telegramAlias": "alias"}`))
				ctx.Request = req

				e.storage.EXPECT().GetUserByAlias(mock.Anything, storeTypes.GetUserByAliasInput{
					Alias: "alias",
				}).Return(storeTypes.User{
					ID:     1,
					Name:   "Timr",
					Alias:  "alias",
					ChatID: 123,
				}, nil)

				e.sender.EXPECT().SendMessageToUserByAlias(
					mock.Anything,
					"alias",
					"Привет! Кое-кто запросил код для входа: `111`").
					Return(assert.AnError) // Ошибка отправки

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
				require.Equal(t, "", w.Body.String())
			},
		},
		{
			name: "FAIL, invalid JSON syntax",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`invalid json`)) // Полностью невалидный JSON
				ctx.Request = req

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
				require.Equal(t, "", w.Body.String())
			},
		},
		{
			name: "FAIL, wrong JSON type for alias",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"telegramAlias": 123}`)) // Число вместо строки
				ctx.Request = req

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
				require.Equal(t, "", w.Body.String())
			},
		},
		{
			name: "OK, user with special characters in alias",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				specialAlias := "user_with_underscore"
				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"telegramAlias": "user_with_underscore"}`))
				ctx.Request = req

				e.storage.EXPECT().GetUserByAlias(mock.Anything, storeTypes.GetUserByAliasInput{
					Alias: specialAlias,
				}).Return(storeTypes.User{
					ID:     1,
					Name:   "Timr",
					Alias:  specialAlias,
					ChatID: 123,
				}, nil)

				e.sender.EXPECT().SendMessageToUserByAlias(
					mock.Anything,
					specialAlias,
					"Привет! Кое-кто запросил код для входа: `111`").
					Return(nil)

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, "", w.Body.String())
			},
		},
		{
			name: "OK, very long alias",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				longAlias := strings.Repeat("a", 100) // Длинный alias
				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"telegramAlias": "`+longAlias+`"}`))
				ctx.Request = req

				e.storage.EXPECT().GetUserByAlias(mock.Anything, storeTypes.GetUserByAliasInput{
					Alias: longAlias,
				}).Return(storeTypes.User{
					ID:     1,
					Name:   "Timr",
					Alias:  longAlias,
					ChatID: 123,
				}, nil)

				e.sender.EXPECT().SendMessageToUserByAlias(
					mock.Anything,
					longAlias,
					"Привет! Кое-кто запросил код для входа: `111`").
					Return(nil)

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, "", w.Body.String())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEnv(t)
			c := newController(e)

			ctx, w := tt.mocks(e)
			c.RequestCode(ctx)

			tt.respCheck(t, ctx, w)
		})
	}
}
func TestController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		mocks       func(e *env) (*gin.Context, *httptest.ResponseRecorder)
		respCheck   func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder)
		checkCookie func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK - successful login",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"alias": "testuser", "otp": "111"}`))
				ctx.Request = req

				e.storage.EXPECT().GetUserByAlias(mock.Anything, storeTypes.GetUserByAliasInput{
					Alias: "testuser",
				}).Return(storeTypes.User{
					ID:     1,
					Name:   "Test User",
					Alias:  "testuser",
					ChatID: 123,
				}, nil)

				e.jwt.On("Sign", web.JWTUserIDKey, domain.UserID(1)).Return("mock-jwt-token", nil)

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, "", w.Body.String())
			},
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				require.Len(t, cookies, 1)
				cookie := cookies[0]
				require.Equal(t, "access-token", cookie.Name)
				require.Equal(t, "mock-jwt-token", cookie.Value)
				require.Equal(t, "/", cookie.Path)
				require.Equal(t, "dengovie.ingress", cookie.Domain)
				require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
				require.False(t, cookie.Secure)
				require.True(t, cookie.HttpOnly)
			},
		},
		{
			name: "FAIL - invalid JSON body",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"alias": "testuser", "otp": "111"`)) // Неполный JSON
				ctx.Request = req

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				require.Len(t, cookies, 0)
			},
		},
		{
			name: "FAIL - invalid OTP",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"alias": "testuser", "otp": "222"}`)) // Неверный OTP
				ctx.Request = req

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
				expectedBody := `{"ErrorReason":"invalid_otp"}`
				require.JSONEq(t, expectedBody, w.Body.String())
			},
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				require.Len(t, cookies, 0)
			},
		},
		{
			name: "FAIL - user not found in storage",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"alias": "nonexistent", "otp": "111"}`))
				ctx.Request = req

				e.storage.EXPECT().GetUserByAlias(mock.Anything, storeTypes.GetUserByAliasInput{
					Alias: "nonexistent",
				}).Return(storeTypes.User{}, sql.ErrNoRows)

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				require.Len(t, cookies, 0)
			},
		},
		{
			name: "FAIL - storage error",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"alias": "testuser", "otp": "111"}`))
				ctx.Request = req

				e.storage.EXPECT().GetUserByAlias(mock.Anything, storeTypes.GetUserByAliasInput{
					Alias: "testuser",
				}).Return(storeTypes.User{}, assert.AnError)

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				require.Len(t, cookies, 0)
			},
		},
		{
			name: "FAIL - JWT signing error",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"alias": "testuser", "otp": "111"}`))
				ctx.Request = req

				e.storage.EXPECT().GetUserByAlias(mock.Anything, storeTypes.GetUserByAliasInput{
					Alias: "testuser",
				}).Return(storeTypes.User{
					ID:     1,
					Name:   "Test User",
					Alias:  "testuser",
					ChatID: 123,
				}, nil)

				e.jwt.On("Sign", web.JWTUserIDKey, domain.UserID(1)).Return("", assert.AnError)

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				require.Len(t, cookies, 0)
			},
		},
		{
			name: "FAIL - wrong data types in JSON",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"alias": 123, "otp": "111"}`)) // alias как число
				ctx.Request = req

				return ctx, w
			},
			respCheck: func(t *testing.T, ctx *gin.Context, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
			checkCookie: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				require.Len(t, cookies, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEnv(t)

			// Создаем контроллер с моком JWT
			c := newController(e)
			ctx, w := tt.mocks(e)
			c.Login(ctx)

			tt.respCheck(t, ctx, w)

			if tt.checkCookie != nil {
				tt.checkCookie(t, w)
			}

			e.jwt.AssertExpectations(t)
		})
	}
}
