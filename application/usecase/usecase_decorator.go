package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/Gabriel-Schiestl/go-clarch/v3/utils"
)

// Classe base abstrata que todos os use cases devem herdar
type BaseUseCase[P, R any] struct {
	processName string
}

// Construtor para a classe base
func NewBaseUseCase[P, R any](processName string) BaseUseCase[P, R] {
	return BaseUseCase[P, R]{
		processName: processName,
	}
}

// Método público que os handlers chamam - equivalente ao seu execute()
func (b *BaseUseCase[P, R]) Execute(ctx context.Context, props P) (R, error) {
	useCaseName := b.processName
	if useCaseName == "" {
		// Fallback para usar reflection se processName não foi definido
		useCaseType := reflect.TypeOf(b)
		useCaseName = useCaseType.String()
		if useCaseType.Kind() == reflect.Ptr {
			useCaseName = useCaseType.Elem().String()
		}
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

	// Chama o método abstrato que deve ser implementado pelos use cases filhos
	result, err := b.OnExecute(ctx, props)

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

// Método abstrato que deve ser implementado pelos use cases filhos
// Equivalente ao seu onExecute()
func (b *BaseUseCase[P, R]) OnExecute(ctx context.Context, props P) (R, error) {
	panic("OnExecute must be implemented by concrete use case")
}

// Interface que define o contrato para os use cases
type ExecutableUseCase[P, R any] interface {
	Execute(ctx context.Context, props P) (R, error)
	OnExecute(ctx context.Context, props P) (R, error)
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