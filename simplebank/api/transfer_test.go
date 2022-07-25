package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/ebaudet/simplebank/db/mock"
	db "github.com/ebaudet/simplebank/db/sqlc"
	"github.com/ebaudet/simplebank/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferAPI(t *testing.T) {
	fromAccount := randomAccount()
	fromAccount.Currency = utils.USD
	toAccount := randomAccount()
	toAccount.Currency = utils.USD
	transfer := randomTransfer(fromAccount, toAccount)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			body: gin.H{
				"from_account_id": transfer.FromAccountID,
				"to_account_id":   transfer.ToAccountID,
				"amount":          transfer.Amount,
				"currency":        fromAccount.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.TransferTxParams{
					FromAccountID: transfer.FromAccountID,
					ToAccountID:   transfer.ToAccountID,
					Amount:        transfer.Amount,
				}

				resultFromAccount := db.Account{
					ID:        fromAccount.ID,
					Owner:     fromAccount.Owner,
					Balance:   fromAccount.Balance - transfer.Amount,
					Currency:  fromAccount.Currency,
					CreatedAt: fromAccount.CreatedAt,
				}
				resultToAccount := db.Account{
					ID:        toAccount.ID,
					Owner:     toAccount.Owner,
					Balance:   toAccount.Balance + transfer.Amount,
					Currency:  toAccount.Currency,
					CreatedAt: toAccount.CreatedAt,
				}
				fromAccountEntry := db.Entry{
					ID:        utils.RandomInt(1, 1000),
					AccountID: fromAccount.ID,
					Amount:    -transfer.Amount,
				}
				toAccountEntry := db.Entry{
					ID:        utils.RandomInt(1, 1000),
					AccountID: toAccount.ID,
					Amount:    transfer.Amount,
				}

				result := db.TransferTxResult{
					Transfer:    transfer,
					FromAccount: resultFromAccount,
					ToAccount:   resultToAccount,
					FromEntry:   fromAccountEntry,
					ToEntry:     toAccountEntry,
				}

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(fromAccount, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
					Times(1).
					Return(toAccount, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchTransfer(t, recorder.Body, transfer)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)
			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomTransfer(fromAccount db.Account, toAccount db.Account) db.Transfer {
	return db.Transfer{
		ID:            utils.RandomInt(1, 1000),
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        utils.RandomMoney(),
	}
}

func requireBodyMatchTransfer(t *testing.T, body *bytes.Buffer, transfer db.Transfer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTransferTxResult db.TransferTxResult
	err = json.Unmarshal(data, &gotTransferTxResult)
	require.NoError(t, err)
	require.Equal(t, transfer, gotTransferTxResult.Transfer)
}
