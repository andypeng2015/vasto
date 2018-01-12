package store

import (
	"github.com/chrislusf/vasto/pb"
	"golang.org/x/net/context"
	"fmt"
)

func (ss *storeServer) DebugStore(ctx context.Context, request *pb.Empty) (*pb.Empty, error) {

	resp := &pb.Empty{}

	ss.debug()

	return resp, nil

}

func (ss *storeServer) debug() () {

	ss.statusInClusterLock.RLock()
	fmt.Println("\n========================================================")
	fmt.Printf("local shards:\n")
	for keyspace, localShards := range ss.statusInCluster {
		fmt.Printf("  * %s\n", keyspace)
		fmt.Printf("    %+v\n", localShards)
	}
	ss.statusInClusterLock.RUnlock()

	fmt.Printf("\nperiodic tasks:\n")
	for _, task := range ss.periodTasks {
		fmt.Printf("  * %v\n", task)
	}

	ss.keyspaceShards.RLock()
	fmt.Printf("\nkeyspace shards:\n")
	for keyspaceName, shards := range ss.keyspaceShards.keyspaceToShards {
		fmt.Printf("  * %v\n", keyspaceName)
		for _, shard := range shards {
			fmt.Printf("    * %v\n", shard.String())
			shard.followProgressLock.Lock()
			for k, v := range shard.followProgress {
				fmt.Printf("      ~ %v.%d @ %d:%d\n", k.serverAdminAddress, k.shardId, v.segment, v.offset)
			}
			shard.followProgressLock.Unlock()
		}
	}
	ss.keyspaceShards.RUnlock()

	fmt.Printf("\ncluster listener event processors:\n")
	ss.clusterListener.Debug("  ")

	fmt.Println()

}
