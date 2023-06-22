package rueidis

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/redis/rueidis"
	"github.com/stretchr/testify/assert"
)

var client rueidis.Client

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	runOptions := &dockertest.RunOptions{
		Repository: "redis",
	}
	resource, err := pool.RunWithOptions(runOptions)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	if err := pool.Retry(func() error {
		client, err = rueidis.NewClient(rueidis.ClientOption{
			InitAddress: []string{fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp"))},
		})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Printf("Could not connect to database: %s", err)
		return 1
	}

	return m.Run()
}

func Test_Redis(t *testing.T) {
	{
		cmd := client.B().Get().Key("key").Build()
		err := client.Do(context.Background(), cmd).Error()
		assert.Error(t, err)
		assert.ErrorIs(t, err, rueidis.Nil)
	}
	{
		cmd := client.B().Set().Key("key").Value("value").Build()
		assert.NoError(t, client.Do(context.Background(), cmd).Error())
	}
	{
		cmd := client.B().Get().Key("key").Build()
		result := client.Do(context.Background(), cmd)
		assert.NoError(t, result.Error())

		str, err := result.ToString()
		assert.NoError(t, err)
		assert.Equal(t, "value", str)
	}
}
