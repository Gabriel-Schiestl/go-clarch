package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/Gabriel-Schiestl/go-clarch/utils"
)

func ExecuteUseCase[R any](ctx context.Context, useCase UseCase[R]) (R, error) {
    useCaseType := reflect.TypeOf(useCase)
    useCaseName := useCaseType.String()
    
    if useCaseType.Kind() == reflect.Ptr {
        useCaseName = useCaseType.Elem().String()
    }

    now := time.Now()

    utils.Logger.Debug().Str("useCase", useCaseName).Msg("Executing use case")

    result, err := useCase.Execute(ctx)

    duration := time.Since(now)

    if err != nil {
        utils.Logger.Error().Err(err).Str("useCase", useCaseName).Str("duration", duration.String()).Msg("Error executing use case")
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

// Execute para casos de uso com parâmetros - agora é uma função genérica
func ExecuteUseCaseWithProps[P, R any](ctx context.Context, useCase UseCaseWithProps[P, R], props P) (R, error) {
    useCaseType := reflect.TypeOf(useCase)
    useCaseName := useCaseType.String()
    
    if useCaseType.Kind() == reflect.Ptr {
        useCaseName = useCaseType.Elem().String()
    }

    propsType := "nil"
    propsValue := reflect.ValueOf(props)
    
    // Verificar se props não é um valor zero
    if !propsValue.IsZero() {
        propsType = propsValue.Type().String()
        
        // Tenta serializar props para logging
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

    result, err := useCase.Execute(ctx, props)

    duration := time.Since(now)

    if err != nil {
        utils.Logger.Error().Err(err).Str("useCase", useCaseName).Str("duration", duration.String()).Msg("Error executing use case")
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

// Função auxiliar para formatar o resultado para logging
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