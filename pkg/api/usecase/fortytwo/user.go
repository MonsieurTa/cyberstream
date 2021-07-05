package fortytwo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type UserInfo struct {
	ID              int       `json:"id"`
	Email           string    `json:"email"`
	Login           string    `json:"login"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	UsualFirstName  string    `json:"usual_first_name"`
	Url             string    `json:"url"`
	Phone           string    `json:"phone"`
	Displayname     string    `json:"displayname"`
	UsualFullName   string    `json:"usual_full_name"`
	ImageUrl        string    `json:"image_url"`
	Staff           bool      `json:"staff?"`
	CorrectionPoint int       `json:"correction_point"`
	PoolMonth       string    `json:"pool_month"`
	PoolYear        string    `json:"pool_year"`
	Location        string    `json:"location"`
	Wallet          int       `json:"wallet"`
	AnonymizeDate   time.Time `json:"anonymize_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Groups          []string  `json:"groups"`
	CursusUsers     []struct {
		ID            int       `json:"id"`
		BeginAt       time.Time `json:"begin_at"`
		EndAt         time.Time `json:"end_at"`
		CursusID      int       `json:"cursus_id"`
		HasCoallition bool      `json:"has_coallition"`
		Grade         string    `json:"grade"`
		Level         float32   `json:"level"`
		BlackholedAt  time.Time `json:"blackholed_at"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Skills        []struct {
			ID    int     `json:"id"`
			Name  string  `json:"name"`
			Level float32 `json:"level"`
		} `json:"skills"`
		User struct {
			ID        int       `json:"id"`
			Login     string    `json:"login"`
			Url       string    `json:"url"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		}
		Cursus struct {
			ID        int       `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			Name      string    `json:"name"`
			Slug      string    `json:"slug"`
		}
	} `json:"cursus_users"`
	// ProjectsUsers   `json:"projects_users"`
	// LanguagesUsers  `json:"languages_users"`
	// Achievements    `json:"achievements"`
	// Titles          `json:"titles"`
	// TitlesUsers     `json:"titles_users"`
	// Partnerships    `json:"partnerships"`
	// Patroned        `json:"patroned"`
	// Patroning       `json:"patroning"`
	// ExpertisesUsers `json:"expertises_users"`
	// Roles           `json:"roles"`
	// Campus          `json:"campus"`
	// CampusUsers     `json:"campus_users"`
}

func (s *Service) GetUserInfo(accessToken string) (*UserInfo, error) {
	baseUrl, err := url.Parse("https://api.intra.42.fr/v2/me")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("access_token", accessToken)

	baseUrl.RawQuery = params.Encode()

	resp, err := http.Get(baseUrl.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var me UserInfo

	err = json.Unmarshal(body, &me)
	if err != nil {
		return nil, err
	}
	return &me, nil
}
