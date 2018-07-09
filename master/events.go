package master

import (
	"github.com/jonas747/dutil/dstate"
)

// Shard rescaling graceful restarts:
// 1. new slave connects
// 2. master sends the new slave EvtSoftStart
// 		This will make the slave start all the shards, but only process the events in the state handler
// 3. the master waits for EvtSoftStartComplete from the new slave
// 4. the master sends the old slave EvtShutdown and the new slave EvtFullStart

// Shard migration using resumes
// 1. new slave connects
// 2. for each shard
// 		a. master sends EvtStopShard to the old slave
//      b. master waits for EvtShardStopped that includes shard info and state info
//      c. master sends EvtResume to the new slave with the new info
// 3. once out of shards the old slave exits and the new slave starts fully

// The event IDs are hardcoded to preserve compatibility between versions
const (
	// Master -> slave
	EvtSoftStart uint32 = 1 // Sent to signal the slave to not start anything other than start updating the state
	EvtFullStart uint32 = 2 // Sent after a soft start event to start up everything other than the state

	// Common, sent by both master and slaves

	// Sent to tell a shard that shard migration is about to happen, either to or from this shard
	// If from this shard to a new one, then responds with the session info needed
	// Otherwise, responds with no data once ready
	EvtShardMigrationStart uint32 = 3

	// Sent to tell the slave to stop a shard, responds with EvtStopShard once all state has been transfered and shard has been stopped
	EvtStopShard uint32 = 4
	// Sent to tell the slave to resume the specified shard, responds with EvtResume once finished
	EvtResume uint32 = 5

	EvtShutdown   uint32 = 6 // Sent to tell a slave to shut down, and immediately stop processing events, responds with the same event once shut down
	EvtGuildState uint32 = 7

	// Slave -> master
	EvtSlaveHello        uint32 = 8
	EvtSoftStartComplete uint32 = 9  // Sent to indicate that all shards has been connected and are waiting for the full start event
	EvtShardStopped      uint32 = 10 // Send by a slave when the shard has been stopped, includes state information for guilds related to that shardd
)

type SlaveHelloData struct {
	Running bool // Wether the slave was already running or not
}

type ShardMigrationStartData struct {
	FromThisSlave bool // If true, were migrating from this shard to a new one, otherwise its the other way around, (from another shard to this one)
	NumShards     int
}

type StopShardData struct {
	Shard int

	SessionID string
	Sequence  int64
}

type ResumeShardData struct {
	Shard int

	SessionID string
	Sequence  int64
}

type GuildStateData struct {
	GuildState *dstate.GuildState
}

var EvtDataMap = map[uint32]func() interface{}{
	EvtSlaveHello:          func() interface{} { return new(SlaveHelloData) },
	EvtShardMigrationStart: func() interface{} { return new(ShardMigrationStartData) },
	EvtStopShard:           func() interface{} { return new(StopShardData) },
	EvtResume:              func() interface{} { return new(ResumeShardData) },
	EvtGuildState:          func() interface{} { return new(GuildStateData) },
}
