package mojangService

import (
	"encoding/json"
	"goa-golang/internal/logger"
	"net/http"
	"time"
)

//MojangServiceInterface define the mojang service interface methods
type MojangServiceInterface interface {
	Profiles(uuid string) (profiles []map[string]interface{}, err error)
}

// billingService handles communication with the mojang repository
type MojangService struct {
	logger logger.Logger
}

// NewMojangService implements the mojang service interface.
func NewMojangService(logger logger.Logger) MojangServiceInterface {
	return &MojangService{logger}
}

// FindByID implements the method to find a mojang model by primary key
func (s *MojangService) Profiles(uuid string) (profiles []map[string]interface{}, err error) {
	s.logger.Infof("uuid: " + uuid)

	url := "https://api.mojang.com/user/profiles/" + uuid + "/names"

	request, _ := http.NewRequest(http.MethodGet, url, nil)

	client := &http.Client{Timeout: 10 * time.Second}

	response, err := client.Do((request))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseJson []map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&responseJson)
	if err != nil {
		return nil, err
	}

	// j, err := json.MarshalIndent(responseJson, "", "\t")
	// if err != nil {
	// 	return nil, err
	// }
	// s.logger.Infof(string(j))

	return responseJson, nil
}
