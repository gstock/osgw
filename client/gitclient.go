package client

import (
"encoding/json"
"log"
"errors"
"net/http"
"github.com/ddliu/go-httpclient"
)


type User struct {
    Login        string   `json:"login,omitempty"`
    Id           int64    `json:"id,omitempty"`
    AvatarUrl    string   `json:"avatar_url,omitempty"`
    Url          string   `json:"url,omitempty"`
    ReposUrl     string   `json:"repos_url,omitempty"`
    Location     string   `json:"location,omitempty"`
}

type Repository struct {
    Id          int64    `json:"id,omitempty"`
    Name        string  `json:"name,omitempty"`
    CreatedAt   string  `json:"created_at,omitempty"`
    UpdatedAt   string  `json:"updated_at,omitempty"`
}

type RepoClient interface { 
    GetUser(string) (User, error)
    GetRepos(string) ([]Repository, error)
} 

type GitClient struct{} 

func (c GitClient) GetUser(username string) (User, error) {
    
    var user User

    if username == "" {        
        return user, errors.New("Username cannot be empty")
    }

    // TODO: don't concatenate the url
    res, err := httpclient.Get("https://api.github.com/users/" + username)

    if err != nil || res.StatusCode != http.StatusOK {
        log.Println("[GitClient > HTTP ERROR] ", err)        
        return user, errors.New("Error obtaining user info from Github")
    }
    

    jsonErr := json.NewDecoder(res.Body).Decode(&user)

    if jsonErr != nil {
        log.Println("[GitClient > JSON ERROR] ", jsonErr)
        return user, errors.New("Error parsing user info from Github")
    }

    return user, nil
}

func (c GitClient) GetRepos(username string) ([]Repository, error) {
    
    var repositories []Repository

    if username == "" {        
        return repositories, errors.New("Username cannot be empty")
    }

    // TODO: don't concatenate the url
    res, err := httpclient.Get("https://api.github.com/users/" + username + "/repos")

    if err != nil || res.StatusCode != http.StatusOK {
        log.Println("[GitClient > HTTP ERROR] ", err)        
        return repositories, errors.New("Error obtaining repositories from Github")
    }

    jsonErr := json.NewDecoder(res.Body).Decode(&repositories)

    if jsonErr != nil {
        log.Println("[GitClient > JSON ERROR] ", jsonErr)
        return repositories, errors.New("Error parsing repositories from Github")
    }

    return repositories, nil
  
}




