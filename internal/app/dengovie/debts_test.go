package dengovie

import (
	"bytes"
	"dengovie/internal/domain"
	debtsTypes "dengovie/internal/service/debts/types"
	storeTypes "dengovie/internal/store/types"
	"dengovie/internal/web"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func TestController_ShareDebt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		prepareMocks    func(t *testing.T, e *env, ctx *gin.Context)
		checkAnswerFunc func(t *testing.T, h *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			checkAnswerFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)
				assert.Empty(t, w.Body)
			},
			prepareMocks: func(t *testing.T, e *env, ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(1))

				jsonRequestValue := Must(json.Marshal(map[string]any{
					"group_id": 2,
					"user_ids": []int64{3, 4},
					"amount":   int64(100),
				}))
				ctx.Request = &http.Request{Header: make(http.Header)}
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonRequestValue))

				e.debtsService.EXPECT().ShareDebt(mock.AnythingOfType("todoCtx"), debtsTypes.ShareDebtInput{
					BuyerID:   domain.UserID(1),
					GroupID:   domain.GroupID(2),
					DebtorIDs: []domain.UserID{3, 4},
					Amount:    int64(100),
				}).Return(nil)

				return
			},
		},
		{
			name: "FAIL, unknown error from debtsService",
			checkAnswerFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
				assert.Empty(t, w.Body)
			},
			prepareMocks: func(t *testing.T, e *env, ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(1))

				jsonRequestValue := Must(json.Marshal(map[string]any{
					"group_id": 2,
					"user_ids": []int64{3, 4},
					"amount":   int64(100),
				}))
				ctx.Request = &http.Request{Header: make(http.Header)}
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonRequestValue))

				e.debtsService.EXPECT().ShareDebt(mock.AnythingOfType("todoCtx"), debtsTypes.ShareDebtInput{
					BuyerID:   domain.UserID(1),
					GroupID:   domain.GroupID(2),
					DebtorIDs: []domain.UserID{3, 4},
					Amount:    int64(100),
				}).Return(assert.AnError)

				return
			},
		},
		{
			name: "FAIL, ErrBuyerNotInGroup",
			checkAnswerFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				assert.NotEmpty(t, w.Body)
				assert.Equal(t, string(Must(json.Marshal(web.APIError{
					ErrorReason: web.BuyerNotInGroup,
				}))), w.Body.String())
			},
			prepareMocks: func(t *testing.T, e *env, ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(1))

				jsonRequestValue := Must(json.Marshal(map[string]any{
					"group_id": 2,
					"user_ids": []int64{3, 4},
					"amount":   int64(100),
				}))
				ctx.Request = &http.Request{Header: make(http.Header)}
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonRequestValue))

				e.debtsService.EXPECT().ShareDebt(mock.AnythingOfType("todoCtx"), debtsTypes.ShareDebtInput{
					BuyerID:   domain.UserID(1),
					GroupID:   domain.GroupID(2),
					DebtorIDs: []domain.UserID{3, 4},
					Amount:    int64(100),
				}).Return(debtsTypes.ErrBuyerNotInGroup)

				return
			},
		},
		{
			name: "FAIL, ErrBuyerNotInGroup",
			checkAnswerFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				assert.NotEmpty(t, w.Body)
				assert.Equal(t, string(Must(json.Marshal(web.APIError{
					ErrorReason: web.DebtorNotInGroup,
				}))), w.Body.String())
			},
			prepareMocks: func(t *testing.T, e *env, ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(1))

				jsonRequestValue := Must(json.Marshal(map[string]any{
					"group_id": 2,
					"user_ids": []int64{3, 4},
					"amount":   int64(100),
				}))
				ctx.Request = &http.Request{Header: make(http.Header)}
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonRequestValue))

				e.debtsService.EXPECT().ShareDebt(mock.AnythingOfType("todoCtx"), debtsTypes.ShareDebtInput{
					BuyerID:   domain.UserID(1),
					GroupID:   domain.GroupID(2),
					DebtorIDs: []domain.UserID{3, 4},
					Amount:    int64(100),
				}).Return(debtsTypes.ErrDebtorNotInGroup)

				return
			},
		},
		{
			name: "FAIL, invalid input",
			checkAnswerFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				assert.Empty(t, w.Body)
			},
			prepareMocks: func(t *testing.T, e *env, ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(1))

				jsonRequestValue := Must(json.Marshal(map[string]any{
					"group_id": "2",
					"user_ids": []int64{3, 4},
					"amount":   int64(100),
				}))
				ctx.Request = &http.Request{Header: make(http.Header)}
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonRequestValue))

				return
			},
		},
		{
			name: "FAIL, no user_id",
			checkAnswerFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				assert.Empty(t, w.Body)
			},
			prepareMocks: func(t *testing.T, e *env, ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, "a")

				return
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			e := newEnv(t)
			tt.prepareMocks(t, e, ctx)
			c := newController(e)

			c.ShareDebt(ctx)

			ctx.Writer.WriteHeaderNow()
			tt.checkAnswerFunc(t, w)
		})
	}
}

func TestController_ListDebts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		checkAnswerFunc func(t *testing.T, h *httptest.ResponseRecorder)
		prepareMocks    func(t *testing.T, e *env, ctx *gin.Context)
	}{
		{
			name: "OK",
			checkAnswerFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)
				assert.NotEmpty(t, w.Body)
				assert.Equal(t, string(Must(json.Marshal(ListDebtsResponseBody{
					Debts: []UserDebt{
						{
							AnotherUserID:   2,
							AnotherUserName: "3",
							Amount:          4,
						},
						{
							AnotherUserID:   5,
							AnotherUserName: "6",
							Amount:          7,
						},
					},
				}))), w.Body.String())
			},
			prepareMocks: func(t *testing.T, e *env, ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(1))

				e.storage.EXPECT().ListUserDebts(mock.AnythingOfType("todoCtx"), storeTypes.ListUserDebtsInput{
					UserID: 1,
				}).Return([]storeTypes.UserDebt{
					{
						AnotherUser: storeTypes.User{
							ID:   2,
							Name: "3",
						},
						Amount: 4,
					},
					{
						AnotherUser: storeTypes.User{
							ID:   5,
							Name: "6",
						},
						Amount: 7,
					},
				}, nil)

				return
			},
		},
		{
			name: "Fail, StatusInternalServerError",
			checkAnswerFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
				assert.Empty(t, w.Body)
			},
			prepareMocks: func(t *testing.T, e *env, ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(1))

				e.storage.EXPECT().ListUserDebts(mock.AnythingOfType("todoCtx"), storeTypes.ListUserDebtsInput{
					UserID: 1,
				}).Return(nil, assert.AnError)

				return
			},
		},
		{
			name: "Fail, StatusInternalServerError",
			checkAnswerFunc: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				assert.Empty(t, w.Body)
			},
			prepareMocks: func(t *testing.T, e *env, ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, "a")
				return
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			e := newEnv(t)
			tt.prepareMocks(t, e, ctx)
			c := newController(e)

			c.ListDebts(ctx)

			ctx.Writer.WriteHeaderNow()
			tt.checkAnswerFunc(t, w)
		})
	}
}

func TestController_PayDebt(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		mocks     func(e *env) (*gin.Context, *httptest.ResponseRecorder)
		setupCtx  func(ctx *gin.Context)
		respCheck func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK - successful pay debt full amount",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"another_user_id": 456, "full": true, "amount": 0}`))
				ctx.Request = req

				e.debtsService.EXPECT().PayDebt(mock.Anything, debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    true,
					Amount:  0,
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
			name: "OK - successful pay debt partial amount",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"another_user_id": 456, "full": false, "amount": 500}`))
				ctx.Request = req

				e.debtsService.EXPECT().PayDebt(mock.Anything, debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    false,
					Amount:  500,
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
					bytes.NewBufferString(`{"another_user_id": 456, "full": true, "amount": 0}`))
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
			name: "FAIL - invalid JSON",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"another_user_id": 456, "full": true`)) // Неполный JSON
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
		{
			name: "FAIL - service error",
			mocks: func(e *env) (*gin.Context, *httptest.ResponseRecorder) {
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)

				req, _ := http.NewRequest("", "",
					bytes.NewBufferString(`{"another_user_id": 456, "full": true, "amount": 0}`))
				ctx.Request = req

				e.debtsService.EXPECT().PayDebt(mock.Anything, debtsTypes.PayDebtInput{
					UserID:  domain.UserID(123),
					PayeeID: domain.UserID(456),
					Full:    true,
					Amount:  0,
				}).Return(assert.AnError)

				return ctx, w
			},
			setupCtx: func(ctx *gin.Context) {
				ctx.Set(domain.UserIDKey, domain.UserID(123))
			},
			respCheck: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEnv(t)
			c := newController(e)

			ctx, w := tt.mocks(e)
			tt.setupCtx(ctx)

			c.PayDebt(ctx)
			tt.respCheck(t, w)
		})
	}
}
