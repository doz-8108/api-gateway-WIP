package utils

import (
	"strconv"

	"github.com/doz-8108/api-gateway/config"
	"github.com/sqids/sqids-go"
)

type (
	Utils struct {
		SqId        *sqids.Sqids
		SqIdInitNum uint64
		EnvVars     config.EnvVars
	}
)

func NewUtils(envVars config.EnvVars) Utils {
	userIdInitNum, err := strconv.ParseUint(envVars.USER_ID_ENCODE_INIT_NUM, 10, 64)
	if err != nil {
		panic(err)
	}

	sqId, _ := sqids.New(sqids.Options{
		Alphabet:  envVars.USER_ID_ENCODE_ALPHB,
		MinLength: 10,
	})

	return Utils{
		SqId:        sqId,
		SqIdInitNum: userIdInitNum,
		EnvVars:     envVars,
	}
}
