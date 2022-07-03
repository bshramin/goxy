package elastic

import (
	"context"
	"fmt"
	"net/http"

	elastic "github.com/olivere/elastic/v7"
	elasticsearch "go.elastic.co/apm/module/apmelasticsearch"
)

func GetConnection(host string, port int) (*elastic.Client, error) {
	connection := fmt.Sprintf("http://%s:%d", host, port)

	elast, err := elastic.NewClient(
		elastic.SetURL(connection),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(&http.Client{
			Transport: elasticsearch.WrapRoundTripper(http.DefaultTransport),
		}),
	)

	return elast, err
}

func HealthCheck(ctx context.Context, client elastic.Client, connection string) error {
	_, _, err := client.Ping(connection).Do(ctx)
	return err
}
