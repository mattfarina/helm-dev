package cmd

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/kube"
	rspb "helm.sh/helm/v3/pkg/release"
)

var settings = cli.New()

func debug(format string, v ...interface{}) {
	if settings.Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

func NewInspectResourceCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "inspect-resource NAME",
		Short: "View the details of a secret resource for a Helm release",
		RunE: func(cmd *cobra.Command, args []string) error {

			kc := kube.New(settings.RESTClientGetter())
			kc.Log = debug

			namespace, present := os.LookupEnv("HELM_NAMESPACE")
			if !present {
				fmt.Println("Namespace is not set")
				return nil
			}

			lazyClient := &lazyClient{
				namespace: namespace,
				clientFn:  kc.Factory.KubernetesClientSet,
			}

			d := driver.NewSecrets(newSecretClient(lazyClient))
			d.Log = debug
			store := storage.Init(d)

			i64, _ := strconv.ParseInt(args[1], 10, 32)
			obj, err := store.Get(args[0], int(i64))
			if err != nil {
				return fmt.Errorf("error getting release: %s", err)
			}
			fmt.Printf("%+v\n", obj)
			return nil
		},
	}

	return rootCmd
}

var b64 = base64.StdEncoding
var magicGzip = []byte{0x1f, 0x8b, 0x08}

func decodeRelease(data string) (*rspb.Release, error) {
	// base64 decode string
	b, err := b64.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// For backwards compatibility with releases that were stored before
	// compression was introduced we skip decompression if the
	// gzip magic header is not found
	if len(b) > 3 && bytes.Equal(b[0:3], magicGzip) {
		r, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		defer r.Close()
		b2, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		b = b2
	}

	var rls rspb.Release
	// unmarshal release object bytes
	if err := json.Unmarshal(b, &rls); err != nil {
		return nil, err
	}
	return &rls, nil
}
