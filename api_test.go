package main

import (
"testing"
"net/http"
"net/http/httptest"
"github.com/gin-gonic/gin"
"osgw/client"
"errors"
"log"
"encoding/json"
)

type MockGitClientInvalid struct{} 

func (c MockGitClientInvalid) GetUser(username string) (client.User, error) {
    return client.User{}, errors.New("Error obtaining user info from Github")
}

func (c MockGitClientInvalid) GetRepos(username string) ([]client.Repository, error) {
    return nil, errors.New("Username cannot be empty")
}

type MockGitClientValid struct{} 

func (c MockGitClientValid) GetUser(username string) (client.User, error) {
    return client.User{Location: "Buenos Aires"}, nil
}

func (c MockGitClientValid) GetRepos(username string) ([]client.Repository, error) {
    repos := []client.Repository{
        client.Repository{Name: "Repo1", CreatedAt: "2016-06-14T11:38:38Z"},
        client.Repository{Name: "Repo2", CreatedAt: "2011-04-10T07:18:20Z"},
        client.Repository{Name: "Repo3", CreatedAt: "2010-10-01T22:12:01Z"},
    }

    return repos, nil
}

type MockGitClientNoRepos struct{} 

func (c MockGitClientNoRepos) GetUser(username string) (client.User, error) {
    return client.User{Location: "Buenos Aires"}, nil
}

func (c MockGitClientNoRepos) GetRepos(username string) ([]client.Repository, error) {
    repos := []client.Repository{        
    }

    return repos, nil
}

type MockWWOClient struct{} 

func (c MockWWOClient) GetWeather(location string, date string) (client.Weather, error) {    
    switch date  {
        case "2016-06-14": return client.Weather{MaxTempC : "25", MinTempC : "10"}, nil
        case "2011-04-10": return client.Weather{MaxTempC : "15", MinTempC : "7"}, nil    
        case "2010-10-01": return client.Weather{MaxTempC : "13", MinTempC : "2"}, nil    
    }
    return client.Weather{MaxTempC : "20", MinTempC : "4"}, nil  
}

type MockWWOClientInvalidDate struct{} 

func (c MockWWOClientInvalidDate) GetWeather(location string, date string) (client.Weather, error) {    
    switch date  {
        case "2016-06-14": return client.Weather{MaxTempC : "25", MinTempC : "10"}, nil
        case "2011-04-10": return client.Weather{MaxTempC : "15", MinTempC : "7"}, nil    
        case "2010-10-01": return client.Weather{}, errors.New("Error obtaining weather from WWO")
    }
    return client.Weather{MaxTempC : "20", MinTempC : "4"}, nil  
}

func TestGetRepoAvgTempNoUser(t *testing.T) {

    responseRecorder := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(responseRecorder)
    c.Params = gin.Params{gin.Param{Key: "username", Value: ""}}
    
    app := &App{repoClient: client.GitClient{}, weatherClient: MockWWOClient{}}

    app.getRepoAvgTemp(c)

    if responseRecorder.Code != http.StatusBadRequest  {
       t.Errorf("Reponse was incorrect, got: %d, want: %d. Error: %s", responseRecorder.Code, http.StatusBadRequest, responseRecorder.Body)
    }
}

func TestGetRepoAvgTemp(t *testing.T) {

    responseRecorder := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(responseRecorder)
    c.Params = gin.Params{gin.Param{Key: "username", Value: "gstock"}}
    
    app := &App{repoClient: client.GitClient{}, weatherClient: MockWWOClient{}}

    app.getRepoAvgTemp(c)

    if responseRecorder.Code != http.StatusOK  {
       t.Errorf("Reponse was incorrect, got: %d, want: %d. Error: %s", responseRecorder.Code, http.StatusOK, responseRecorder.Body)
    }
}

func TestInvalidUser(t *testing.T) {

    responseRecorder := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(responseRecorder)
    c.Params = gin.Params{gin.Param{Key: "username", Value: "johndoe"}}
    
    app := &App{repoClient: MockGitClientInvalid{}, weatherClient: MockWWOClient{}}

    app.getRepoAvgTemp(c)

    if responseRecorder.Code != http.StatusServiceUnavailable  {
       t.Errorf("Reponse was incorrect, got: %d, want: %d. Error: %s", responseRecorder.Code, http.StatusServiceUnavailable, responseRecorder.Body)
    }
}

func TestValidUser(t *testing.T) {

    responseRecorder := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(responseRecorder)
    c.Params = gin.Params{gin.Param{Key: "username", Value: "johndoe"}}
    
    app := &App{repoClient: MockGitClientValid{}, weatherClient: MockWWOClient{}}

    app.getRepoAvgTemp(c)

    if responseRecorder.Code != http.StatusOK  {
       t.Errorf("Reponse was incorrect, got: %d, want: %d. Error: %s", responseRecorder.Code, http.StatusServiceUnavailable, responseRecorder.Body)
    }

    var response RepoAvgTempResponse
    jsonErr := json.NewDecoder(responseRecorder.Body).Decode(&response)
    
    if response.Count != 3 || response.AvgTemp != 12.0 {
        t.Errorf("Reponse was incorrect, got: %d and %f, want: %d and %f. Error: %s", response.Count, response.AvgTemp, 3, 12.0, jsonErr)
    }
}

func TestNoRepos(t *testing.T) {

    responseRecorder := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(responseRecorder)
    c.Params = gin.Params{gin.Param{Key: "username", Value: "gstock"}}
    
    app := &App{repoClient: MockGitClientNoRepos{}, weatherClient: MockWWOClient{}}

    app.getRepoAvgTemp(c)

    if responseRecorder.Code != http.StatusOK  {
       t.Errorf("Reponse was incorrect, got: %d, want: %d. Error: %s", responseRecorder.Code, http.StatusServiceUnavailable, responseRecorder.Body)
    }

    log.Println(responseRecorder.Body)
}

func TestValidUserInvalidDate(t *testing.T) {

    responseRecorder := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(responseRecorder)
    c.Params = gin.Params{gin.Param{Key: "username", Value: "johndoe"}}
    
    app := &App{repoClient: MockGitClientValid{}, weatherClient: MockWWOClientInvalidDate{}}

    app.getRepoAvgTemp(c)

    if responseRecorder.Code != http.StatusOK  {
       t.Errorf("Reponse was incorrect, got: %d, want: %d. Error: %s", responseRecorder.Code, http.StatusServiceUnavailable, responseRecorder.Body)
    }

    var response RepoAvgTempResponse
    jsonErr := json.NewDecoder(responseRecorder.Body).Decode(&response)
    
    if response.Count != 3 || response.AvgTemp != 14.25 {
        t.Errorf("Reponse was incorrect, got: %d and %f, want: %d and %f. Error: %s", response.Count, response.AvgTemp, 3, 14.25, jsonErr)
    }
}


