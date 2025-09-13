package dengovie

import (
	"bytes"
	"dengovie/internal/domain"
	usersTypes "dengovie/internal/service/users/types"
	storeTypes "dengovie/internal/store/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestController_GetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		mocks     func(e *env) (*gin.Context, *httptest.ResponseRecorder)
		setupCtx  func(ctx *gin.Context)
		respCheck func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK - successful get me",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				e.storage.EXPECT().GetUserByID(ctx, storeTypes.GetUserByIDInput{
					UserID: domain.UserID(123),
				}).Return(storeTypes.User{
					ID:    123,
					Name:  "John Doe",
					Alias: "johndoe",
				}, nil)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				expectedBody := `{
					"user_id": 123,
					"name": "John Doe",
					"alias": "johndoe"
				}`
				require.JSONEq(t, expectedBody, w.Body.String())
			},
		},
		{
			name: "FAIL - userID not found in context",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)
				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				// Не устанавливаем userID
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "FAIL - storage error",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				e.storage.EXPECT().GetUserByID(ctx, storeTypes.GetUserByIDInput{
					UserID: domain.UserID(123),
				}).Return(storeTypes.User{}, assert.AnError)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "FAIL - invalid userID type in context",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)
				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				// Устанавливаем неверный тип userID
				ctx.Set(domain.UserIDKey, "invalid-type")
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "OK - user with empty name and alias",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				e.storage.EXPECT().GetUserByID(ctx, storeTypes.GetUserByIDInput{
					UserID: domain.UserID(123),
				}).Return(storeTypes.User{
					ID:    123,
					Name:  "",
					Alias: "",
				}, nil)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				expectedBody := `{
					"user_id": 123,
					"name": "",
					"alias": ""
				}`
				require.JSONEq(t, expectedBody, w.Body.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEnv(t)
			c := newController(e)

			ctx, w := tt.mocks(e)
			tt.setupCtx(ctx)

			c.GetMe(ctx)
			tt.respCheck(t, w)
		})
	}
}

func TestController_UpdateName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		mocks     func(e *env) (*gin.Context, *httptest.ResponseRecorder)
		setupCtx  func(ctx *gin.Context)
		respCheck func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK - successful update name",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"new_name": "John Updated"}`))
				ctx.Request = req

				e.storage.EXPECT().UpdateUserName(ctx, storeTypes.UpdateUserNameInput{
					UserID:  domain.UserID(123),
					NewName: "John Updated",
				}).Return(nil)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "FAIL - userID not found in context",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"new_name": "John Updated"}`))
				ctx.Request = req

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				// Не устанавливаем userID
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "FAIL - storage error",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"new_name": "John Updated"}`))
				ctx.Request = req

				e.storage.EXPECT().UpdateUserName(ctx, storeTypes.UpdateUserNameInput{
					UserID:  domain.UserID(123),
					NewName: "John Updated",
				}).Return(assert.AnError)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "FAIL - very long name",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				longName := strings.Repeat("a", 1000) // Очень длинное имя
				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"new_name": "`+longName+`"}`))
				ctx.Request = req

				e.storage.EXPECT().UpdateUserName(ctx, storeTypes.UpdateUserNameInput{
					UserID:  domain.UserID(123),
					NewName: longName,
				}).Return(nil)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "FAIL - wrong data type for new_name",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"new_name": 123}`)) // Число вместо строки
				ctx.Request = req

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEnv(t)
			c := newController(e)

			ctx, w := tt.mocks(e)
			tt.setupCtx(ctx)

			c.UpdateName(ctx)
			tt.respCheck(t, w)
		})
	}
}

func TestController_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		mocks     func(e *env) (*gin.Context, *httptest.ResponseRecorder)
		setupCtx  func(ctx *gin.Context)
		respCheck func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK - successful delete user",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				e.usersService.EXPECT().CheckAndDeleteUser(ctx, usersTypes.CheckAndDeleteUserInput{
					UserID: domain.UserID(123),
				}).Return(nil)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "FAIL - userID not found in context",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)
				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				// Не устанавливаем userID
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "FAIL - service error",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				e.usersService.EXPECT().CheckAndDeleteUser(ctx, usersTypes.CheckAndDeleteUserInput{
					UserID: domain.UserID(123),
				}).Return(assert.AnError)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "FAIL - invalid userID type in context",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)
				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				// Устанавливаем неверный тип userID
				ctx.Set(domain.UserIDKey, "invalid-type")
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "OK - delete user with different HTTP method (DELETE)",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				// Используем DELETE метод
				req, _ := http.NewRequest("DELETE", "/user/delete", nil)
				ctx.Request = req

				e.usersService.EXPECT().CheckAndDeleteUser(ctx, usersTypes.CheckAndDeleteUserInput{
					UserID: domain.UserID(123),
				}).Return(nil)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEnv(t)
			c := newController(e)

			ctx, w := tt.mocks(e)
			tt.setupCtx(ctx)

			c.DeleteUser(ctx)
			tt.respCheck(t, w)
		})
	}
}
