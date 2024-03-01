package api

import(
	"github.com/gin-gonic/gin"
	"solution/pkg/models"
)


func GetAllCountries(c *gin.Context){
	region := c.Query("region")
	countries, err := models.GetAllCountries(region)
	if err != nil{
		c.JSON(400, gin.H{"error":"Формат входного запроса не соответствует формату либо переданы неверные значения"})
	} else{
		c.JSON(200,countries)
	}
}
func GetCountryByid(c *gin.Context){
	alpha2 := c.Param("alpha2")
	country, err := models.GetCountryByid(alpha2)
	
	if err != nil{
		c.JSON(404, gin.H{"error":"Страна с указанным кодом не найдена."})
	} else{
		response := models.CountryResponse{
			Name:   country.Name,
			Alpha2: country.Alpha2,
			Alpha3: country.Alpha3,
			Region: country.Region,
		}
		c.JSON(200,response)
	}

}

  