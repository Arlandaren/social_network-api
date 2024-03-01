package api

import(
	"github.com/gin-gonic/gin"
	"solution/pkg/models"
)


func GetAllCountries(c *gin.Context){
	region := c.Query("region")
	countries, err := models.GetAllCountries(region)
	
	if err != nil{
		c.JSON(500, gin.H{"error":err.Error()})
	} else{
		c.JSON(200,countries)
	}
}
func GetCountryByid(c *gin.Context){
	alpha2 := c.Param("alpha2")
	country, err := models.GetCountryByid(alpha2)
	response := models.CountryResponse{
        Name:   country.Name,
        Alpha2: country.Alpha2,
        Alpha3: country.Alpha3,
        Region: country.Region,
    }

	if err != nil{
		c.JSON(500, gin.H{"error":err.Error()})
	} else{
		c.JSON(200,response)
	}

}

  