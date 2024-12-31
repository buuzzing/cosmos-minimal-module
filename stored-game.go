package checkers

import (
	"fmt"

	"cosmossdk.io/errors"
	"github.com/buzzing/checkers/rules"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetBlackAddress 返回游戏对局中黑方的地址
func (storedGame *StoredGame) GetBlackAddress() (black sdk.AccAddress, err error) {
	black, errBlack := sdk.AccAddressFromBech32(storedGame.Black)
	return black, errors.Wrapf(errBlack, ErrInvalidBlack.Error(), storedGame.Black)
}

// GetRedAddress 返回游戏对局中红方的地址
func (storedGame *StoredGame) GetRedAddress() (red sdk.AccAddress, err error) {
	red, errRed := sdk.AccAddressFromBech32(storedGame.Red)
	return red, errors.Wrapf(errRed, ErrInvalidRed.Error(), storedGame.Red)
}

// ParseGame 解析（反序列化）游戏对局
func (storedGame *StoredGame) ParseGame() (game *rules.Game, err error) {
	board, errBoard := rules.Parse(storedGame.Board)
	if errBoard != nil {
		return nil, errors.Wrapf(errBoard, ErrGameNotParseable.Error())
	}
	board.Turn = rules.StringPieces[storedGame.Turn].Player
	if board.Turn.Color == "" {
		return nil, errors.Wrapf(fmt.Errorf("turn: %s", storedGame.Turn), ErrGameNotParseable.Error())
	}
	return board, nil
}

// Validate 验证游戏对局
func (storedGame *StoredGame) Validate() (err error) {
	_, err = storedGame.GetBlackAddress()
	if err != nil {
		return err
	}
	_, err = storedGame.GetRedAddress()
	if err != nil {
		return err
	}
	_, err = storedGame.ParseGame()
	return err
}
