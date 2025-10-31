package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	"github.com/buzzing/checkers"
)

// SendRecordPacket 发送 RecordPacketData 数据包
// 触发跨链发送的入口，通过指定的 sourcePort 和 sourceChannel 发送数据包
func (k *Keeper) SendRecordPacket(
	ctx sdk.Context,
	packetData checkers.RecordPacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	// 获取通道的 Capability，确保模块有权限使用该通道
	// OCAP 模型，参见 keeper/keeper.go 中的介绍
	channelCap, ok := k.ScopedKeeper().GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, errorsmod.Wrapf(channeltypes.ErrChannelCapabilityNotFound,
			"module does not own channel capability for port %s channel %s", sourcePort, sourceChannel)
	}

	packetDataBytes, err := checkers.ModuleCdc.MarshalJSON(&packetData)
	if err != nil {
		return 0, errorsmod.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %v", err)
	}

	return k.ibcKeeperFn().ChannelKeeper.SendPacket(ctx, channelCap, sourcePort, sourceChannel,
		timeoutHeight, timeoutTimestamp, packetDataBytes)
}

// OnRecvRecordPacket 处理接收到的 RecordPacketData 数据包
// 自定义处理逻辑的函数
func (k *Keeper) OnRecvRecordPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	data checkers.RecordPacketData,
) (packetAck checkers.RecordPacketAck, err error) {
	// TODO: 自定义处理逻辑

	return packetAck, err
}

// OnAcknowledgementRecordPacket 处理 RecordPacketData 数据包的确认
// 自定义处理逻辑的函数
func (k *Keeper) OnAcknowledgementRecordPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	data checkers.RecordPacketData,
	acknowledgement channeltypes.Acknowledgement,
) error {
	// TODO: 自定义处理逻辑

	return nil
}

// OnTimeoutRecordPacket 处理 RecordPacketData 数据包的超时
// 自定义处理逻辑的函数
func (k *Keeper) OnTimeoutRecordPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	data checkers.RecordPacketData,
) error {
	// TODO: 自定义处理逻辑

	return nil
}
