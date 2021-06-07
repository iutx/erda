package clusters

import (
	"github.com/erda-project/erda/apistructs"
	"os"
)

func (c *Clusters) InitClusters(req apistructs.ClusterInitRequest) (uint64, error) {
	erdaClusterName := os.Getenv(string(apistructs.DICE_CLUSTER_NAME))
	if erdaClusterName == "" {
		// TODO
	}

	return 0, nil
}
