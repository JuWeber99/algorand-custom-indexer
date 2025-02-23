package encoding

import (
	"github.com/algorand/indexer/types"

	sdk "github.com/algorand/go-algorand-sdk/v2/types"
)

type blockHeader struct {
	sdk.BlockHeader
	BranchOverride      sdk.Digest `codec:"prev"`
	FeeSinkOverride     sdk.Digest `codec:"fees"`
	RewardsPoolOverride sdk.Digest `codec:"rwd"`
}

type assetParams struct {
	sdk.AssetParams
	UnitNameBytes    []byte     `codec:"un64"`
	AssetNameBytes   []byte     `codec:"an64"`
	URLBytes         []byte     `codec:"au64"`
	ManagerOverride  sdk.Digest `codec:"m"`
	ReserveOverride  sdk.Digest `codec:"r"`
	FreezeOverride   sdk.Digest `codec:"f"`
	ClawbackOverride sdk.Digest `codec:"c"`
}

type transaction struct {
	sdk.Transaction
	SenderOverride           sdk.Digest   `codec:"snd"`
	RekeyToOverride          sdk.Digest   `codec:"rekey"`
	ReceiverOverride         sdk.Digest   `codec:"rcv"`
	CloseRemainderToOverride sdk.Digest   `codec:"close"`
	AssetParamsOverride      assetParams  `codec:"apar"`
	AssetSenderOverride      sdk.Digest   `codec:"asnd"`
	AssetReceiverOverride    sdk.Digest   `codec:"arcv"`
	AssetCloseToOverride     sdk.Digest   `codec:"aclose"`
	FreezeAccountOverride    sdk.Digest   `codec:"fadd"`
	AccountsOverride         []sdk.Digest `codec:"apat"`
}

type valueDelta struct {
	sdk.ValueDelta
	BytesOverride []byte `codec:"bs"`
}

type byteArray struct {
	data string
}

func (ba byteArray) MarshalText() ([]byte, error) {
	return []byte(Base64([]byte(ba.data))), nil
}

func (ba *byteArray) UnmarshalText(text []byte) error {
	baNew, err := decodeBase64(string(text))
	if err != nil {
		return err
	}

	*ba = byteArray{string(baNew)}
	return nil
}

type stateDelta map[byteArray]valueDelta

type evalDelta struct {
	sdk.EvalDelta
	GlobalDeltaOverride stateDelta            `codec:"gd"`
	LocalDeltasOverride map[uint64]stateDelta `codec:"ld"`
	LogsOverride        [][]byte              `codec:"lg"`
	InnerTxnsOverride   []signedTxnWithAD     `codec:"itx"`
}

type signedTxnWithAD struct {
	sdk.SignedTxnWithAD
	TxnOverride       transaction `codec:"txn"`
	AuthAddrOverride  sdk.Digest  `codec:"sgnr"`
	EvalDeltaOverride evalDelta   `codec:"dt"`
}

type trimmedAccountData struct {
	baseAccountData
	AuthAddrOverride   sdk.Digest `codec:"spend"`
	MicroAlgos         uint64     `codec:"algo"`
	RewardsBase        uint64     `codec:"ebase"`
	RewardedMicroAlgos uint64     `codec:"ern"`
}

type tealValue struct {
	sdk.TealValue
	BytesOverride []byte `codec:"tb"`
}

type tealKeyValue map[byteArray]tealValue

type appLocalState struct {
	sdk.AppLocalState
	KeyValueOverride tealKeyValue `codec:"tkv"`
}

type appParams struct {
	sdk.AppParams
	GlobalStateOverride tealKeyValue `codec:"gs"`
}

type specialAddresses struct {
	types.SpecialAddresses
	FeeSinkOverride     sdk.Digest `codec:"FeeSink"`
	RewardsPoolOverride sdk.Digest `codec:"RewardsPool"`
}

type baseOnlineAccountData struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	VoteID          sdk.OneTimeSignatureVerifier `codec:"vote"`
	SelectionID     sdk.VRFVerifier              `codec:"sel"`
	StateProofID    sdk.Commitment               `codec:"stprf"`
	VoteFirstValid  sdk.Round                    `codec:"voteFst"`
	VoteLastValid   sdk.Round                    `codec:"voteLst"`
	VoteKeyDilution uint64                       `codec:"voteKD"`
}

type baseAccountData struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Status              sdk.Status      `codec:"onl"`
	AuthAddr            sdk.Digest      `codec:"spend"`
	TotalAppSchema      sdk.StateSchema `codec:"tsch"`
	TotalExtraAppPages  uint32          `codec:"teap"`
	TotalAssetParams    uint64          `codec:"tasp"`
	TotalAssets         uint64          `codec:"tas"`
	TotalAppParams      uint64          `codec:"tapp"`
	TotalAppLocalStates uint64          `codec:"tapl"`
	TotalBoxes          uint64          `codec:"tbx"`
	TotalBoxBytes       uint64          `codec:"tbxb"`

	baseOnlineAccountData
}
