package timeout

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestTimeoutMiddleware(t *testing.T) {

	router := gin.Default()
	router.Use(Timeout(time.Second))
	router.GET("/", func(c *gin.Context) {
		select {
		case <-time.After(time.Millisecond):
			c.String(200, "success")
		case <-c.Request.Context().Done():
			c.String(408, "timed out")
		}
	})

	srv := httptest.NewServer(router)
	defer srv.Close()

	res, err := http.Get(srv.URL)
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)

	message, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, "success", string(message))
}

func TestTimeoutMiddlewareFailure(t *testing.T) {

	router := gin.Default()
	router.Use(Timeout(time.Second))
	router.GET("/", func(c *gin.Context) {
		select {
		case <-time.After(time.Minute):
			c.String(200, "success")
		case <-c.Request.Context().Done():
			c.String(408, "timed out")
		}
	})

	srv := httptest.NewServer(router)
	defer srv.Close()

	res, err := http.Get(srv.URL)
	require.NoError(t, err)
	require.Equal(t, 408, res.StatusCode)

	message, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)
	require.Equal(t, "timed out", string(message))
}
