package module

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"github.com/buzzing/checkers"
)

var (
	// AppModule 需要实现 IBCModule 接口
	_ porttypes.IBCModule = AppModule{}
)

// OnChanOpenInit 在通道开启初始化时被调用
// 模块可以在这里执行自定义逻辑，例如是否接受该通道
func (am AppModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portId string,
	channelId string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	// 检查 portId 是否是模块绑定的 Port
	if portId != checkers.PortId {
		return "", errorsmod.Wrapf(porttypes.ErrInvalidPort,
			"invalid port: %s, expected %s", portId, checkers.PortId)
	}
	// 检查 version 是否匹配
	if version != checkers.Version {
		return "", errorsmod.Wrapf(checkers.ErrInvalidVersion,
			"invalid version: %s, expected %s", version, checkers.Version)
	}

	// 声明模块拥有这个通道的能力
	if err := am.keeper.ClaimCapability(ctx, channelCap, host.ChannelCapabilityPath(portId, channelId)); err != nil {
		return "", err
	}

	return version, nil
}

// OnChanOpenTry 在通道开启尝试时被调用
// 通道的另一端使用
func (am AppModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portId,
	channelId string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	// 检查 portId 是否是模块绑定的 Port
	if portId != checkers.PortId {
		return "", errorsmod.Wrapf(porttypes.ErrInvalidPort,
			"invalid port: %s, expected %s", portId, checkers.PortId)
	}
	// 检查 version 是否匹配
	if counterpartyVersion != checkers.Version {
		return "", errorsmod.Wrapf(checkers.ErrInvalidVersion,
			"invalid version: %s, expected %s", counterpartyVersion, checkers.Version)
	}

	// 模块可能已经存储了这个 Chan 的 Capability（例如在 OnChanOpenInit 中）
	// 因此这里使用 AuthenticateCapability 来检查
	if !am.keeper.AuthenticateCapability(ctx, channelCap, host.ChannelCapabilityPath(portId, channelId)) {
		// 如果没有存储，则认领这个 Capability
		if err := am.keeper.ClaimCapability(ctx, channelCap, host.ChannelCapabilityPath(portId, channelId)); err != nil {
			return "", err
		}
	}

	return counterpartyVersion, nil
}

// OnChanOpenAck 在通道开启确认时被调用
func (am AppModule) OnChanOpenAck(
	ctx sdk.Context,
	portId string,
	channelId string,
	_,
	counterpartyVersion string,
) error {
	if counterpartyVersion != checkers.Version {
		return errorsmod.Wrapf(checkers.ErrInvalidVersion,
			"invalid counterparty version: %s, expected %s", counterpartyVersion, checkers.Version)
	}
	return nil
}

// OnChanOpenConfirm 在通道开启最终确认时被调用
func (am AppModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portId string,
	channelId string,
) error {
	return nil
}

// OnChanCloseInit 在通道关闭初始化时被调用
func (am AppModule) OnChanCloseInit(
	ctx sdk.Context,
	portId string,
	channelId string,
) error {
	// 模块不允许主动关闭通道
	return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "module does not allow channel closure")
}

// OnChanCloseConfirm 在通道关闭确认时被调用
func (am AppModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portId string,
	channelId string,
) error {
	return nil
}

// OneRecvPacket 在接收数据包时被调用
// 当模块收到一个 IBC 数据包时，此函数被调用
func (am AppModule) OnRecvPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	var ack channeltypes.Acknowledgement

	var modulePacketData checkers.RecordPacketData
	if err := modulePacketData.Unmarshal(modulePacket.GetData()); err != nil {
		return channeltypes.NewErrorAcknowledgement(
			errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %s", err.Error()),
		)
	}

	packetAck, err := am.keeper.OnRecvRecordPacket(ctx, modulePacket, modulePacketData)
	if err != nil {
		ack = channeltypes.NewErrorAcknowledgement(err)
	} else {
		// 成功处理后，创建 ACK 消息并编码
		packetAckBytes, err := checkers.ModuleCdc.MarshalJSON(&packetAck)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(
				errorsmod.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal packet acknowledgement: %s", err.Error()),
			)
		}
		ack = channeltypes.NewResultAcknowledgement(packetAckBytes)
	}

	return ack
}

// OnAcknowledgementPacket 在数据包确认时被调用
func (am AppModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := checkers.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet acknowledgement: %v", err)
	}

	var modulePacketData checkers.RecordPacketData
	if err := modulePacketData.Unmarshal(modulePacket.GetData()); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %v", err)
	}

	var eventType string

	err := am.keeper.OnAcknowledgementRecordPacket(ctx, modulePacket, modulePacketData, ack)
	if err != nil {
		return err
	}
	eventType = checkers.EventTypeIbcRecordPacket

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(sdk.AttributeKeyModule, checkers.ModuleName),
			sdk.NewAttribute(checkers.AttributeKeyAck, fmt.Sprintf("%v", ack)),
		),
	)

	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				eventType,
				sdk.NewAttribute(checkers.AttributeKeyAckSuccess, string(resp.Result)),
			),
		)
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				eventType,
				sdk.NewAttribute(checkers.AttributeKeyAckError, resp.Error),
			),
		)
	}

	return nil
}

// OnTimeoutPacket 在数据包超时时被调用
func (am AppModule) OnTimeoutPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	var modulePacketData checkers.RecordPacketData
	if err := modulePacketData.Unmarshal(modulePacket.GetData()); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet data: %v", err)
	}

	return am.keeper.OnTimeoutRecordPacket(ctx, modulePacket, modulePacketData)
}
