package api

import(
	"github.com/gin-gonic/gin"
	"solution/models"
)

func GetAllCountries(c *gin.Context){
	region := c.Query("region")
	countries, err := models.GetAllCountries(region)
	if err != nil{
		c.JSON(500, gin.H{"error":err.Error()})
	} else{
		c.JSON(200,gin.H{"countries":countries})
	}
}
func GetCountryByid(c *gin.Context){
	alpha2 := c.Param("alpha2")
	country, err := models.GetCountryByid(alpha2)
	if err != nil{
		c.JSON(500, gin.H{"error":err.Error()})
	} else{
		c.JSON(200,gin.H{"countries":country})
	}

}

  