package elastic

import (
	"context"
	"fmt"
	"net/http"

	elastic "github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	elasticsearch "go.elastic.co/apm/module/apmelasticsearch"
)

func GetConnection(host string, port int) *elastic.Client {
	connection := fmt.Sprintf("http://%s:%d", host, port)

	elast, err := elastic.NewClient(
		elastic.SetURL(connection),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(&http.Client{
			Transport: elasticsearch.WrapRoundTripper(http.DefaultTransport),
		}),
	)
	if err != nil {
		logrus.Error("connect to elastic failed")
		panic(err)
	}

	logrus.Infof("successfully connected to elastic : %s", connection)
	return elast
}

func HealthCheck(ctx context.Context, client elastic.Client, connection string) error {

	_, _, err := client.Ping(connection).Do(ctx)
	return err
}
