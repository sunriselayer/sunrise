/**
 * @file abci_test.go
 * @brief DAモジュールのABCIインターフェース関数のテスト
 *
 * このファイルでは、DAモジュールのEndBlocker関数のテストを実装しています。
 * 主に以下の機能をテストします：
 * - 期限切れのREJECTEDステータスデータの削除
 * - 期限切れのVERIFIEDステータスデータの削除
 * - CHALLENGE_PERIODステータスのデータの処理
 * - CHALLENGINGステータスのデータの処理
 * - スラッシュエポックの処理
 */

package keeper_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sunriselayer/sunrise/x/da/keeper"
	"github.com/sunriselayer/sunrise/x/da/types"
)

func TestEndBlocker(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(fixture)
		expectError bool
	}{
		{
			name: "削除期限切れのREJECTEDステータスデータ",
			setup: func(f fixture) {
				// パラメータの設定
				params, err := f.keeper.Params.Get(f.ctx)
				require.NoError(t, err)
				params.RejectedRemovalPeriod = time.Hour * 24
				err = f.keeper.Params.Set(f.ctx, params)
				require.NoError(t, err)

				// 期限切れのREJECTEDデータを作成
				oldTime := sdk.UnwrapSDKContext(f.ctx).BlockTime().Add(-time.Hour * 48)
				data := types.PublishedData{
					MetadataUri: "test-uri-1",
					Status:      types.Status_STATUS_REJECTED,
					Timestamp:   oldTime,
				}
				err = f.keeper.SetPublishedData(f.ctx, data)
				require.NoError(t, err)

				// 期限内のREJECTEDデータも作成
				newTime := sdk.UnwrapSDKContext(f.ctx).BlockTime().Add(-time.Hour)
				data2 := types.PublishedData{
					MetadataUri: "test-uri-2",
					Status:      types.Status_STATUS_REJECTED,
					Timestamp:   newTime,
				}
				err = f.keeper.SetPublishedData(f.ctx, data2)
				require.NoError(t, err)

				// GetSpecificStatusDataBeforeTimeのモック
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_REJECTED, gomock.Any()).
					Return([]types.PublishedData{data}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_VERIFIED, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusData(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGING, gomock.Any()).
					Return([]types.PublishedData{}, nil)

				// DeletePublishedDataのモック
				f.mocks.DAKeeper.EXPECT().
					DeletePublishedData(gomock.Any(), data).
					Return(nil)
			},
		},
		{
			name: "削除期限切れのVERIFIEDステータスデータ",
			setup: func(f fixture) {
				// パラメータの設定
				params, err := f.keeper.Params.Get(f.ctx)
				require.NoError(t, err)
				params.VerifiedRemovalPeriod = time.Hour * 24
				err = f.keeper.Params.Set(f.ctx, params)
				require.NoError(t, err)

				// 期限切れのVERIFIEDデータを作成
				oldTime := sdk.UnwrapSDKContext(f.ctx).BlockTime().Add(-time.Hour * 48)
				data := types.PublishedData{
					MetadataUri: "test-uri-1",
					Status:      types.Status_STATUS_VERIFIED,
					Timestamp:   oldTime,
				}
				err = f.keeper.SetPublishedData(f.ctx, data)
				require.NoError(t, err)

				// GetSpecificStatusDataBeforeTimeのモック
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_REJECTED, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_VERIFIED, gomock.Any()).
					Return([]types.PublishedData{data}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusData(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGING, gomock.Any()).
					Return([]types.PublishedData{}, nil)

				// DeletePublishedDataのモック
				f.mocks.DAKeeper.EXPECT().
					DeletePublishedData(gomock.Any(), data).
					Return(nil)
			},
		},
		{
			name: "CHALLENGE_PERIODステータスのデータをCHALLENGINGに変更",
			setup: func(f fixture) {
				// パラメータの設定
				params, err := f.keeper.Params.Get(f.ctx)
				require.NoError(t, err)
				params.ChallengeThreshold = "0.5"
				err = f.keeper.Params.Set(f.ctx, params)
				require.NoError(t, err)

				// CHALLENGE_PERIODデータを作成
				data := types.PublishedData{
					MetadataUri:       "test-uri-1",
					Status:            types.Status_STATUS_CHALLENGE_PERIOD,
					ShardDoubleHashes: []string{"hash1", "hash2", "hash3", "hash4"},
				}
				err = f.keeper.SetPublishedData(f.ctx, data)
				require.NoError(t, err)

				// 無効性の作成
				invalidities := []types.Invalidity{
					{
						MetadataUri: "test-uri-1",
						Sender:      "sender1",
						Indices:     []int64{0, 1},
					},
					{
						MetadataUri: "test-uri-1",
						Sender:      "sender2",
						Indices:     []int64{2},
					},
				}

				// GetSpecificStatusDataBeforeTimeのモック
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_REJECTED, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_VERIFIED, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusData(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD).
					Return([]types.PublishedData{data}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetInvalidities(gomock.Any(), "test-uri-1").
					Return(invalidities, nil)
				f.mocks.DAKeeper.EXPECT().
					SetPublishedData(gomock.Any(), gomock.Any()).
					Return(nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGING, gomock.Any()).
					Return([]types.PublishedData{}, nil)
			},
		},
		{
			name: "期限切れのCHALLENGE_PERIODステータスをVERIFIEDに変更",
			setup: func(f fixture) {
				// パラメータの設定
				params, err := f.keeper.Params.Get(f.ctx)
				require.NoError(t, err)
				params.ChallengePeriod = time.Hour * 24
				err = f.keeper.Params.Set(f.ctx, params)
				require.NoError(t, err)

				// 期限切れのCHALLENGE_PERIODデータを作成
				oldTime := sdk.UnwrapSDKContext(f.ctx).BlockTime().Add(-time.Hour * 48)
				data := types.PublishedData{
					MetadataUri:           "test-uri-1",
					Status:                types.Status_STATUS_CHALLENGE_PERIOD,
					Timestamp:             oldTime,
					Publisher:             "publisher1",
					PublishDataCollateral: sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
				}
				err = f.keeper.SetPublishedData(f.ctx, data)
				require.NoError(t, err)

				// GetSpecificStatusDataBeforeTimeのモック
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_REJECTED, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_VERIFIED, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusData(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD, gomock.Any()).
					Return([]types.PublishedData{data}, nil)
				f.mocks.DAKeeper.EXPECT().
					SetPublishedData(gomock.Any(), gomock.Any()).
					Return(nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGING, gomock.Any()).
					Return([]types.PublishedData{}, nil)

				// BankKeeperのモック
				publisherAddr, _ := sdk.AccAddressFromBech32("publisher1")
				f.mocks.BankKeeper.EXPECT().
					SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, publisherAddr, data.PublishDataCollateral).
					Return(nil)
			},
		},
		{
			name: "CHALLENGINGステータスのデータを処理（REJECTEDに変更）",
			setup: func(f fixture) {
				// パラメータの設定
				params, err := f.keeper.Params.Get(f.ctx)
				require.NoError(t, err)
				params.ProofPeriod = time.Hour * 24
				params.ReplicationFactor = "3.0"
				err = f.keeper.Params.Set(f.ctx, params)
				require.NoError(t, err)

				// CHALLENGINGデータを作成
				oldTime := sdk.UnwrapSDKContext(f.ctx).BlockTime().Add(-time.Hour * 48)
				data := types.PublishedData{
					MetadataUri:                "test-uri-1",
					Status:                     types.Status_STATUS_CHALLENGING,
					Timestamp:                  oldTime,
					ShardDoubleHashes:          []string{"hash1", "hash2", "hash3", "hash4"},
					ParityShardCount:           1,
					Publisher:                  "publisher1",
					PublishDataCollateral:      sdk.NewCoins(sdk.NewInt64Coin("stake", 100)),
					SubmitInvalidityCollateral: sdk.NewCoins(sdk.NewInt64Coin("stake", 10)),
				}
				err = f.keeper.SetPublishedData(f.ctx, data)
				require.NoError(t, err)

				// 無効性の作成
				invalidities := []types.Invalidity{
					{
						MetadataUri: "test-uri-1",
						Sender:      "challenger1",
						Indices:     []int64{0, 1},
					},
				}

				// プルーフの作成
				proofs := []types.Proof{
					{
						MetadataUri: "test-uri-1",
						Sender:      "validator1",
						Indices:     []int64{2},
					},
				}

				// GetSpecificStatusDataBeforeTimeのモック
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_REJECTED, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_VERIFIED, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusData(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGING, gomock.Any()).
					Return([]types.PublishedData{data}, nil)

				// ValidatorsPowerStoreIteratorのモック
				f.mocks.StakingKeeper.EXPECT().
					ValidatorsPowerStoreIterator(gomock.Any()).
					Return(nil, nil)

				// GetProofsのモック
				f.mocks.DAKeeper.EXPECT().
					GetProofs(gomock.Any(), "test-uri-1").
					Return(proofs, nil)

				// GetZkpThresholdのモック
				f.mocks.DAKeeper.EXPECT().
					GetZkpThreshold(gomock.Any(), uint64(4)).
					Return(uint64(2), nil)

				// GetInvaliditiesのモック
				f.mocks.DAKeeper.EXPECT().
					GetInvalidities(gomock.Any(), "test-uri-1").
					Return(invalidities, nil)

				// SetPublishedDataのモック
				f.mocks.DAKeeper.EXPECT().
					SetPublishedData(gomock.Any(), gomock.Any()).
					Return(nil)

				// BankKeeperのモック
				challengerAddr, _ := sdk.AccAddressFromBech32("challenger1")
				reward := sdk.NewCoins(sdk.NewInt64Coin("stake", 100))
				f.mocks.BankKeeper.EXPECT().
					SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, challengerAddr, data.SubmitInvalidityCollateral.Add(reward...)).
					Return(nil)

				// SetChallengeCounterのモック
				f.mocks.DAKeeper.EXPECT().
					GetChallengeCounter(gomock.Any()).
					Return(uint64(0))
				f.mocks.DAKeeper.EXPECT().
					SetChallengeCounter(gomock.Any(), uint64(1)).
					Return(nil)

				// DeleteProofのモック
				f.mocks.DAKeeper.EXPECT().
					addressCodec.
					StringToBytes("validator1").
					Return([]byte("validator1"), nil)
				f.mocks.DAKeeper.EXPECT().
					DeleteProof(gomock.Any(), "test-uri-1", []byte("validator1")).
					Return(nil)

				// DeleteInvalidityのモック
				f.mocks.DAKeeper.EXPECT().
					addressCodec.
					StringToBytes("challenger1").
					Return([]byte("challenger1"), nil)
				f.mocks.DAKeeper.EXPECT().
					DeleteInvalidity(gomock.Any(), "test-uri-1", []byte("challenger1")).
					Return(nil)
			},
		},
		{
			name: "スラッシュエポックの処理",
			setup: func(f fixture) {
				// パラメータの設定
				params, err := f.keeper.Params.Get(f.ctx)
				require.NoError(t, err)
				params.SlashEpoch = 100
				err = f.keeper.Params.Set(f.ctx, params)
				require.NoError(t, err)

				// ブロック高さの設定
				ctx := sdk.UnwrapSDKContext(f.ctx)
				ctx = ctx.WithBlockHeight(100)
				f.ctx = ctx

				// GetSpecificStatusDataBeforeTimeのモック
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_REJECTED, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_VERIFIED, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusData(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGE_PERIOD, gomock.Any()).
					Return([]types.PublishedData{}, nil)
				f.mocks.DAKeeper.EXPECT().
					GetSpecificStatusDataBeforeTime(gomock.Any(), types.Status_STATUS_CHALLENGING, gomock.Any()).
					Return([]types.PublishedData{}, nil)

				// HandleSlashEpochのモック
				f.mocks.DAKeeper.EXPECT().
					HandleSlashEpoch(gomock.Any()).
					Return()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := initFixture(t)
			ctx := f.ctx
			k := f.keeper
			mocks := f.mocks

			suite := fixture{
				ctx:    ctx,
				keeper: k,
				mocks:  mocks,
			}
			tt.setup(suite)

			err := k.EndBlocker(ctx)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCheckCorrectInvalidity(t *testing.T) {
	tests := []struct {
		name             string
		invalidity       types.Invalidity
		safeShardIndices []int64
		expected         bool
	}{
		{
			name: "すべての無効性インデックスが安全でない場合",
			invalidity: types.Invalidity{
				Indices: []int64{0, 1, 2},
			},
			safeShardIndices: []int64{3, 4, 5},
			expected:         true,
		},
		{
			name: "一部の無効性インデックスが安全な場合",
			invalidity: types.Invalidity{
				Indices: []int64{0, 1, 3},
			},
			safeShardIndices: []int64{3, 4, 5},
			expected:         false,
		},
		{
			name: "すべての無効性インデックスが安全な場合",
			invalidity: types.Invalidity{
				Indices: []int64{3, 4, 5},
			},
			safeShardIndices: []int64{3, 4, 5},
			expected:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := keeper.CheckCorrectInvalidity(tt.invalidity, tt.safeShardIndices)
			require.Equal(t, tt.expected, result)
		})
	}
}

// fixtureはテスト用の構造体
type fixture struct {
	ctx    context.Context
	keeper keeper.Keeper
	mocks  DAMocks
}

// DAMocksはDAモジュールのモック
type DAMocks struct {
	DAKeeper      *MockDAKeeper
	BankKeeper    *MockBankKeeper
	StakingKeeper *MockStakingKeeper
}

// initFixtureはテスト用のfixtureを初期化する
func initFixture(t *testing.T) fixture {
	ctrl := gomock.NewController(t)
	daKeeper := NewMockDAKeeper(ctrl)
	bankKeeper := NewMockBankKeeper(ctrl)
	stakingKeeper := NewMockStakingKeeper(ctrl)

	k := keeper.Keeper{
		Params:        daKeeper.Params,
		BankKeeper:    bankKeeper,
		StakingKeeper: stakingKeeper,
		addressCodec:  daKeeper.addressCodec,
	}

	ctx := sdk.Context{}.WithBlockTime(time.Now())

	return fixture{
		ctx:    ctx,
		keeper: k,
		mocks: DAMocks{
			DAKeeper:      daKeeper,
			BankKeeper:    bankKeeper,
			StakingKeeper: stakingKeeper,
		},
	}
}

// MockDAKeeperはDAKeeperのモック
type MockDAKeeper struct {
	Params       types.ParamsI
	addressCodec *MockAddressCodec
}

func NewMockDAKeeper(ctrl *gomock.Controller) *MockDAKeeper {
	return &MockDAKeeper{
		Params:       NewMockParamsI(ctrl),
		addressCodec: NewMockAddressCodec(ctrl),
	}
}

// MockParamsIはParamsIのモック
type MockParamsI struct {
	*gomock.Controller
}

func NewMockParamsI(ctrl *gomock.Controller) *MockParamsI {
	return &MockParamsI{Controller: ctrl}
}

func (m *MockParamsI) Get(ctx context.Context) (types.Params, error) {
	return types.DefaultParams(), nil
}

func (m *MockParamsI) Set(ctx context.Context, params types.Params) error {
	return nil
}

// MockAddressCodecはAddressCodecのモック
type MockAddressCodec struct {
	*gomock.Controller
}

func NewMockAddressCodec(ctrl *gomock.Controller) *MockAddressCodec {
	return &MockAddressCodec{Controller: ctrl}
}

func (m *MockAddressCodec) StringToBytes(address string) ([]byte, error) {
	return []byte(address), nil
}

// MockBankKeeperはBankKeeperのモック
type MockBankKeeper struct {
	*gomock.Controller
}

func NewMockBankKeeper(ctrl *gomock.Controller) *MockBankKeeper {
	return &MockBankKeeper{Controller: ctrl}
}

// MockStakingKeeperはStakingKeeperのモック
type MockStakingKeeper struct {
	*gomock.Controller
}

func NewMockStakingKeeper(ctrl *gomock.Controller) *MockStakingKeeper {
	return &MockStakingKeeper{Controller: ctrl}
}
