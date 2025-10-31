package checkers

// IBC 相关事件类型
const (
	EventTypeTimeout         = "timeout"
	EventTypeIbcRecordPacket = "record_packet"

	AttributeKeyAckSuccess = "ack_success"
	AttributeKeyAck        = "ack"
	AttributeKeyAckError   = "ack_error"
)
