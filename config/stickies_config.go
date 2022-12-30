package config

import (
	"os"
	"strconv"
)

type StickiesConfig struct {
	WalnutStickerSetName string
	WalnutStickerUserID  int64
}

func ReadStickiesConfig() (StickiesConfig, error) {
	setName := os.Getenv("WALNUT_STICKER_SET_NAME")
	id, err := strconv.ParseInt(os.Getenv("WALNUT_STICKER_USER_ID"), 10, 64)
	if err != nil {
		return StickiesConfig{}, err
	}

	return StickiesConfig{
		WalnutStickerSetName: setName,
		WalnutStickerUserID:  id,
	}, nil
}
