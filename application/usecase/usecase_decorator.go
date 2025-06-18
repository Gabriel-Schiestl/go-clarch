package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Gabriel-Schiestl/go-clarch/utils"
)

// Decorador para casos de uso sem parâmetros
type UseCaseDecorator[R any] struct {
    useCase UseCase[R]
}

// Decorador para casos de uso com parâmetros
type UseCaseWithPropsDecorator[P, R any] struct {
    useCase UseCaseWithProps[P, R]
}

// Criação do decorador para casos de uso sem parâmetros
func NewUseCaseDecorator[R any](useCase UseCase[R]) UseCaseDecorator[R] {
    return UseCaseDecorator[R]{
        useCase: useCase,
    }
}

// Criação do decorador para casos de uso com parâmetros
func NewUseCaseWithPropsDecorator[P, R any](useCase UseCaseWithProps[P, R]) UseCaseWithPropsDecorator[P, R] {
    return UseCaseWithPropsDecorator[P, R]{
        useCase: useCase,
    }
}

// Execute para casos de uso sem parâmetros
func (d UseCaseDecorator[R]) Execute(ctx context.Context) (R, error) {
    useCaseType := reflect.TypeOf(d.useCase)
    useCaseName := useCaseType.String()
    
    if useCaseType.Kind() == reflect.Ptr {
        useCaseName = useCaseType.Elem().String()
    }

    utils.Logger.Debug().Str("useCase", useCaseName).Msg("Executing use case")

    result, err := d.useCase.Execute(ctx)
    if err != nil {
        utils.Logger.Error().Err(err).Str("useCase", useCaseName).Msg("Error executing use case")
        var zero R
        return zero, err
    }

    resultValue := formatResult(result)

    utils.Logger.Info().
        Str("useCase", useCaseName).
        Str("result", resultValue).
        Msg("Successfully executed use case")
        
    return result, nil
}

// Execute para casos de uso com parâmetros
func (d UseCaseWithPropsDecorator[P, R]) Execute(ctx context.Context, props P) (R, error) {
    useCaseType := reflect.TypeOf(d.useCase)
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

    result, err := d.useCase.Execute(ctx, props)
    if err != nil {
        utils.Logger.Error().Err(err).Str("useCase", useCaseName).Msg("Error executing use case")
        var zero R
        return zero, err
    }

    resultValue := formatResult(result)

    utils.Logger.Info().
        Str("useCase", useCaseName).
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