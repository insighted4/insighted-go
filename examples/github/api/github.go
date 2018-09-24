package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/insighted4/insighted-go/kit/extensions"
)

func (s service) getUserHTTPHandler(c *gin.Context) {
	username := c.Param("id")

	body, err := s.GetUser(context.Context(c), &GetUserRequest{Name: username})
	if err != nil {
		extensions.AbortWithStatusJSON(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, body)
}

func (s service) GetUser(ctx context.Context, req *GetUserRequest) (*User, error) {
	u, r, err := s.client.Users.Get(ctx, req.Name)
	if err != nil && r.StatusCode != http.StatusNotFound {
		return nil, err
	}

	return &User{
		Id:        int64(*u.ID),
		Name:      *u.Name,
		Followers: int64(*u.Followers),
		Following: int64(*u.Following),
	}, nil
}
