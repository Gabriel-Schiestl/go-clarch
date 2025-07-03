package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/Gabriel-Schiestl/go-clarch/v3/utils"
)

type abstractUseCase[P, R any] interface {
    onExecute(ctx context.Context, props P) (R, error)
}

type BaseUseCase[P, R any] struct {
    processName string
    executor    abstractUseCase[P, R]
}

func NewBaseUseCase[P, R any](processName string, executor abstractUseCase[P, R]) BaseUseCase[P, R] {
    return BaseUseCase[P, R]{
        processName: processName,
        executor:    executor,
    }
}

func (b *BaseUseCase[P, R]) Execute(ctx context.Context, props P) (R, error) {
    useCaseName := b.processName
    if useCaseName == "" {
        useCaseType := reflect.TypeOf(b.executor)
        useCaseName = useCaseType.String()
        if useCaseType.Kind() == reflect.Ptr {
            useCaseName = useCaseType.Elem().String()
        }
    }

    propsType := "nil"
    propsValue := reflect.ValueOf(props)

    if !propsValue.IsZero() {
        propsType = propsValue.Type().String()

        propsJson, err := json.Marshal(props)
        if err == nil {
            utils.Logger.Debug().
                Str("useCase", useCaseName).
                Str("propsType", propsType).
                RawJSON("props", propsJson).
                Msg("Executing use case with props")
        } else {
            utils.Logger.Debug().
                Str("useCase", useCaseName).
                Str("propsType", propsType).
                Msg("Executing use case with props")
        }
    } else {
        utils.Logger.Debug().
            Str("useCase", useCaseName).
            Msg("Executing use case with nil props")
    }

    now := time.Now()

    result, err := b.executor.onExecute(ctx, props)

    duration := time.Since(now)

    if err != nil {
        utils.Logger.Error().
            Err(err).
            Str("useCase", useCaseName).
            Str("duration", duration.String()).
            Msg("Error executing use case")
        var zero R
        return zero, err
    }

    resultValue := formatResult(result)

    utils.Logger.Info().
        Str("useCase", useCaseName).
        Str("duration", duration.String()).
        Str("result", resultValue).
        Msg("Successfully executed use case")

    return result, nil
}

func formatResult(result any) string {
	if result == nil {
		return "nil"
	}

	resultValue := "nil"
	resultReflect := reflect.ValueOf(result)

	isNil := false

	switch resultReflect.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface, reflect.Chan, reflect.Func:
		if !resultReflect.IsValid() || resultReflect.IsZero() {
			isNil = true
		} else {
			isNil = resultReflect.IsNil()
		}
	}

	if !isNil {
		jsonData, jsonErr := json.Marshal(result)
		if jsonErr == nil {
			resultValue = string(jsonData)
		} else {
			resultValue = fmt.Sprintf("%+v", result)
		}
	}

	return resultValue
}