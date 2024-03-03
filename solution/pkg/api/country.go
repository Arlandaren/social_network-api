package api

import (
	"solution/pkg/models"

	"github.com/gin-gonic/gin"
)

func GetAllCountries(c *gin.Context) {
	regions := c.QueryArray("region")

	countries := make([]models.CountryResponse, 0)

	if len(regions) == 0{
		regions = append(regions, "none")
	}
	for _, region := range regions {
		if region == ""{
			c.JSON(400, gin.H{"reason": "Формат входного запроса не соответствует формату либо переданы неверные значения"})
			return
		}
		countriesForRegion, err := models.GetAllCountries(region)
		if err != nil {
			c.JSON(400, gin.H{"reason": "Формат входного запроса не соответствует формату либо переданы неверные значения"})
			return
		}
		countries = append(countries, countriesForRegion...)
	}
	

	c.JSON(200, countries)
}
func GetCountryByid(c *gin.Context) {
	alpha2 := c.Param("alpha2")
	country, err := models.GetCountryByid(alpha2)

	if err != nil {
		c.JSON(404, gin.H{"reason": "Страна с указанным кодом не найдена."})
	} else {
		response := models.CountryResponse{
			Name:   country.Name,
			Alpha2: country.Alpha2,
			Alpha3: country.Alpha3,
			Region: country.Region,
		}
		c.JSON(200, response)
	}

}
