// OSGW API app
//
//     Schemes: http
//     Host: localhost:8080
//     BasePath: /
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
//"encoding/json"
"log"
"osgw/client"
"net/http"
"time"
"strconv"
"math"
"github.com/gin-contrib/cors"
"github.com/gin-gonic/gin"
)

// HTTP status code 200 and ticker price
// swagger:response swaggRepoAvgTempResponse
type swaggRepoAvgTempResponse struct {
    // in:body
    Body RepoAvgTempResponse
}

// HTTP status code 503
// swagger:response swaggErrResponse
type swaggErrResponse struct {
    // in:body
    Body struct {
         // HTTP status code 503 -  Service Unavailable
         Error int `json:"error"`
         // Detailed error message
         Message string `json:"message"`
    }
}

type RepoAvgTempResponse struct {
    Count   int   `json:"count,omitempty"`
    AvgTemp float64 `json:"avg_temp,omitempty"`
}

type App struct {
    repoClient client.RepoClient
    weatherClient client.WeatherClient
}


func main() {
    app := &App{repoClient: client.GitClient{}, weatherClient: client.WWOClient{}}

    r := gin.Default()
    r.Use(cors.Default()) // I'm doing this to allow Swagger "try" functionality locally
    r.GET("/api/osgw/:username", app.getRepoAvgTemp)
    r.Run() // listen and serve on 0.0.0.0:8080
}

// swagger:operation GET /api/osgw/{username} repository temperature
// ---
// summary: Returns repository count and temperature
// description: If connection to the providers fails Services Unavailable (503) will be returned.
// parameters:
// - name: username
//   in: path
//   description: username to search
//   type: string
//   required: true
// responses:
//   "200":
//     "$ref": "#/responses/swaggRepoAvgTempResponse"
//   "503":
//     "$ref": "#/responses/swaggErrResponse"

func (a *App) getRepoAvgTemp(c *gin.Context) {
    username := c.Param("username")

    if username == "" {
        c.JSON(http.StatusBadRequest, gin.H{
                "success" : false,
                "error" : "You need to specify the username", 
            })
        return
    }

    user, err := a.repoClient.GetUser(username)
    
    if err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
                "success" : false,
                "error" : err.Error(),
            })
        return
    }

    repositories, err := a.repoClient.GetRepos(username)

    if err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
                "success" : false,
                "error" : err.Error(),
            })
        return
    }
    
    sumTemp := 0.0
    count := 0

    for _, repository := range repositories {
        repoCreationDate, err := parseDate(repository.CreatedAt)

        if err == nil {

            weather, err := a.weatherClient.GetWeather(user.Location, repoCreationDate.Format("2006-01-02"))    

            if err == nil {
                max, errMax := strconv.ParseFloat(weather.MaxTempC, 32)
                min, errMin := strconv.ParseFloat(weather.MinTempC, 32)

                if errMax == nil && errMin == nil {
                    count++
                    avg := (max + min) / 2
                    log.Println("  *", repository.Name, repoCreationDate.Format("2006-01-02"), avg)    
                    sumTemp += avg 
                } else {
                    log.Println("[Error Parsing Weather]", errMax, errMin)
                }
                
            } else {
                log.Println("[Error Obtaining Weather]", err)
            }
        } else {
            log.Println("[Error Parsing Date]", err)
        }
        
    }

    avgTemp := 0.0
    if count > 0 {
        avgTemp = math.Round(sumTemp / float64(count) * 100) / 100
    }

    response := RepoAvgTempResponse{Count : len(repositories), AvgTemp : avgTemp}

    c.JSON(http.StatusOK, response)
}

func parseDate(date string) (time.Time, error) {
    return time.Parse(time.RFC3339, date)
}
