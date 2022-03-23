package checkpoint

import (
	"github.com/tkeel-io/rule-util/stream/types"
)

type Barrier struct {
	checkpointId int64
}

func NewCheckpointContext(checkpointId int64) types.FunctionSnapshotContext {
	return &Barrier{
		checkpointId,
	}
}

func (this *Barrier) GetCheckpointId() types.CheckpointID {
	return types.CheckpointID(this.checkpointId)
}

func (this *Barrier) GetCheckpointTimestamp() types.CheckpointTimestamp {
	return types.CheckpointTimestamp(this.checkpointId)
}
