package weather

import (
	"context"

	"github.com/gerps2/desafio-cloud-run/features/weather/getWeatherByCep"
	sharedErrors "github.com/gerps2/desafio-cloud-run/shared/errors"
	httpShared "github.com/gerps2/desafio-cloud-run/shared/http"
	"github.com/gerps2/desafio-cloud-run/shared/logger"

	"github.com/gin-gonic/gin"
)

type WeatherController struct {
	getWeatherByCepUseCase getWeatherByCep.GetWeatherByCepUseCaseInterface
	logger                 logger.Logger
}

func NewWeatherController(getWeatherByCepUseCase getWeatherByCep.GetWeatherByCepUseCaseInterface, logger logger.Logger) *WeatherController {
	return &WeatherController{
		getWeatherByCepUseCase: getWeatherByCepUseCase,
		logger:                 logger,
	}
}

func (wc *WeatherController) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.GET("/weather/:cep", wc.GetWeatherByCep)
	}
}

func (wc *WeatherController) GetWeatherByCep(c *gin.Context) {
	wc.logger.Info("GetWeatherByCep endpoint called")

	cepParam := c.Param("cep")
	if cepParam == "" {
		wc.logger.Error("CEP parameter is required")
		httpShared.RespondWithValidationError(c, "CEP parameter is required", []string{"CEP parameter must be provided in the URL path"})
		return
	}

	input := getWeatherByCep.GetWeatherByCepInput{
		CepString: cepParam,
	}

	result, err := wc.getWeatherByCepUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if c.Request.Context().Err() == context.DeadlineExceeded {
			wc.logger.Error("Request timeout exceeded for CEP: %s", cepParam)
			return
		}
		wc.logger.Error("Error executing GetWeatherByCep use case: %v", err)

		if apiErr, ok := err.(*sharedErrors.APIError); ok {
			httpShared.RespondWithAPIError(c, apiErr)
		} else {
			httpShared.RespondWithInternalError(c, "Failed to get weather data", []string{err.Error()})
		}
		return
	}

	wc.logger.Info("Weather data retrieved successfully for CEP: %s", cepParam)
	httpShared.RespondWithSuccess(c, result, "Weather data retrieved successfully")
}
