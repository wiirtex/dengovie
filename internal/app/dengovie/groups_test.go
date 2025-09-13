package dengovie

import (
	"dengovie/internal/domain"
	storeTypes "dengovie/internal/store/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestController_ListUserGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		mocks     func(e *env) (*gin.Context, *httptest.ResponseRecorder)
		setupCtx  func(ctx *gin.Context)
		respCheck func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK - successful list user groups",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				e.storage.EXPECT().ListUserGroups(mock.Anything, storeTypes.ListUserGroupsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.Group{
					{
						ID:   1,
						Name: "Family Group",
					},
					{
						ID:   2,
						Name: "Friends Group",
					},
				}, nil)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				expectedBody := `[
					{"ID": 1, "Name": "Family Group"},
					{"ID": 2, "Name": "Friends Group"}
				]`
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

				e.storage.EXPECT().ListUserGroups(mock.Anything, storeTypes.ListUserGroupsInput{
					UserID: domain.UserID(123),
				}).Return(nil, assert.AnError)

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
			name: "OK - empty groups list",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				e.storage.EXPECT().ListUserGroups(mock.Anything, storeTypes.ListUserGroupsInput{
					UserID: domain.UserID(123),
				}).Return([]storeTypes.Group{}, nil)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, "[]", w.Body.String())
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEnv(t)
			c := newController(e)

			ctx, w := tt.mocks(e)
			tt.setupCtx(ctx)

			c.ListUserGroups(ctx)
			tt.respCheck(t, w)
		})
	}
}

func TestController_ListUsersInGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		mocks        func(e *env) (*gin.Context, *httptest.ResponseRecorder)
		setupCtx     func(ctx *gin.Context)
		setupRequest func(req *http.Request) // Установка параметров пути
		respCheck    func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK - successful list users in group",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				e.storage.EXPECT().ListUsersInGroup(mock.Anything, storeTypes.ListUsersInGroupInput{
					GroupID: domain.GroupID(1),
				}).Return([]storeTypes.User{
					{
						ID:    123,
						Name:  "John Doe",
						Alias: "johndoe",
					},
					{
						ID:    456,
						Name:  "Jane Smith",
						Alias: "janesmith",
					},
				}, nil)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(999)) // UserID для проверки авторизации
				params := gin.Params{gin.Param{Key: "groupID", Value: "1"}}
				ctx.Params = params
			},
			setupRequest: func(req *http.Request) {
				gin.SetMode(gin.TestMode)
				ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
				ctx.Request = req
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				expectedBody := `[
					{"ID": 123, "Name": "John Doe", "Alias": "johndoe", "ChatID": 0},
					{"ID": 456, "Name": "Jane Smith", "Alias": "janesmith", "ChatID": 0}
				]`
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
				params := gin.Params{gin.Param{Key: "groupID", Value: "1"}}
				ctx.Params = params
			},
			setupRequest: func(req *http.Request) {
				gin.SetMode(gin.TestMode)
				ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
				ctx.Request = req
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "FAIL - invalid groupID parameter (not a number)",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)
				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(999))
				params := gin.Params{gin.Param{Key: "groupID", Value: "invalid"}}
				ctx.Params = params
			},
			setupRequest: func(req *http.Request) {
				gin.SetMode(gin.TestMode)
				ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
				ctx.Request = req
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "FAIL - missing groupID parameter",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)
				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(999))
				params := gin.Params{}
				ctx.Params = params
			},
			setupRequest: func(req *http.Request) {
				gin.SetMode(gin.TestMode)
				ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
				ctx.Request = req
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

				e.storage.EXPECT().ListUsersInGroup(mock.Anything, storeTypes.ListUsersInGroupInput{
					GroupID: domain.GroupID(1),
				}).Return(nil, assert.AnError)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(999))
				params := gin.Params{gin.Param{Key: "groupID", Value: "1"}}
				ctx.Params = params
			},
			setupRequest: func(req *http.Request) {
				gin.SetMode(gin.TestMode)
				ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
				ctx.Request = req
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "OK - empty users list in group",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				e.storage.EXPECT().ListUsersInGroup(mock.Anything, storeTypes.ListUsersInGroupInput{
					GroupID: domain.GroupID(1),
				}).Return([]storeTypes.User{}, nil)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(999))
				params := gin.Params{gin.Param{Key: "groupID", Value: "1"}}
				ctx.Params = params
			},
			setupRequest: func(req *http.Request) {
				gin.SetMode(gin.TestMode)
				ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
				ctx.Request = req
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				require.Equal(t, "[]", w.Body.String())
			},
		},
		{
			name: "FAIL - very large groupID",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)
				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(999))
				params := gin.Params{gin.Param{Key: "groupID", Value: "9999999999999999999"}}
				ctx.Params = params
			},
			setupRequest: func(req *http.Request) {
				gin.SetMode(gin.TestMode)
				ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
				ctx.Request = req
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

			// Настраиваем параметры запроса
			req, _ := http.NewRequest("GET", "/groups/1/users", nil)
			if tt.setupRequest != nil {
				tt.setupRequest(req)
			}
			ctx.Request = req

			c.ListUsersInGroup(ctx)
			tt.respCheck(t, w)
		})
	}
}

// Вспомогательная функция для создания контекста с параметрами пути
func createTestContextWithParams(w *httptest.ResponseRecorder, params gin.Params) *gin.Context {
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = params
	return ctx
}
