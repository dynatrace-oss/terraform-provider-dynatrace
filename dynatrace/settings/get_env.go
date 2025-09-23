package settings

import (
	"context"
	"os"
	"strconv"
)

func GetIntEnv(name string, def, min, max int) int {
	sValue := os.Getenv(name)
	if len(sValue) == 0 {
		return def
	}
	iValue, err := strconv.Atoi(sValue)
	if err != nil {
		return def
	}
	if iValue < min || iValue > max {
		return def
	}
	return iValue
}

func GetCtxStringValue(ctx context.Context, name string) string {
	if ctxValue := ctx.Value(name); ctxValue != nil {
		if ctxSValue, ok := ctxValue.(string); ok {
			return ctxSValue
		}
	}
	return ""
}

func GetIntEnvCtx(ctx context.Context, name string, def, min, max int) int {
	sValue := GetCtxStringValue(ctx, name)
	if len(sValue) == 0 {
		sValue = os.Getenv(name)
	}
	if len(sValue) == 0 {
		return def
	}
	iValue, err := strconv.Atoi(sValue)
	if err != nil {
		return def
	}
	if iValue < min || iValue > max {
		return def
	}
	return iValue
}
