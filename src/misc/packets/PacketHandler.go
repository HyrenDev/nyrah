package packets

import (
	"errors"
	"fmt"
	"net/hyren/nyrah/minecraft/chat"
	"net/hyren/nyrah/minecraft/protocol"
	"net/hyren/nyrah/minecraft/protocol/codecs"
	"net/hyren/nyrah/minecraft/protocol/packet"
	"reflect"
	"strings"

	ProxyApp "net/hyren/nyrah/applications"
	ProxyConnector "net/hyren/nyrah/misc/connector"
	Constants "net/hyren/nyrah/misc/constants"
	Config "net/hyren/nyrah/misc/utils"
	User "net/hyren/nyrah/users"
)

func HandlePackets(connection *protocol.Connection, holder packet.Holder) error {
	switch connection.State {
	case protocol.Handshake:
		{
			handshake, ok := holder.(packet.Handshake)

			if !ok {
				return errors.New(fmt.Sprintf("expected handshake, received: %s", reflect.TypeOf(holder)))
			}

			connection.Protocol = uint16(handshake.ProtocolVersion)
			connection.State = protocol.State(uint8(handshake.NextState))

			handshake.NextState = 2
			handshake.ServerAddress = codecs.String(string(handshake.ServerAddress) + "%ABC%" + strings.Split(connection.Handle.RemoteAddr().String(), ":")[0])

			fmt.Println("Received ping request from:", handshake.ServerAddress)

			connection.PacketQueue[0] = handshake

			return nil
		}
	case protocol.Status:
		{
			_, ok := holder.(packet.StatusRequest)

			var online = Config.GetOnlinePlayers()
			var maxPlayers = Config.GetMaxPlayers()

			if ok {
				response := packet.StatusResponse{}

				response.Status.Version.Name = "1.8.9"
				response.Status.Version.Protocol = 47
				response.Status.Players.Max = maxPlayers
				response.Status.Players.Online = online
				response.Status.Description = Config.GetMOTD()
				response.Status.ModInfo.Type = "FML"
				response.Status.ModInfo.ModList = []string{}

				favicon, err := Config.GetFavicon()

				if err != nil {
					fmt.Println("Cannot find favicon")
				} else {
					response.Status.Favicon = favicon
				}

				_, _ = connection.Write(response)

				return nil
			}

			statusPing, ok := holder.(packet.StatusPing)

			if ok {
				response := packet.StatusPong{}

				response.Payload = statusPing.Payload

				_, _ = connection.Write(response)

				return nil
			}
		}
	case protocol.Login:
		{
			loginStart, ok := holder.(packet.LoginStart)

			if ok {
				name := string(loginStart.Username)

				fmt.Println(fmt.Sprintf(
					"Conexão recebida de [%s/%s]",
					name,
					connection.GetRemoteAddr(),
				))

				if Config.IsMaintenanceModeEnabled() == true && !User.IsHelperOrHigher(
					name,
				) {
					disconnectBecauseMaintenanceModeIsEnabled(
						connection,
					)
					return nil
				}

				proxies, err := ProxyApp.FetchAvailableProxiesNames()

				if err != nil {
					disconnectBecauseNotHaveProxyToSend(
						connection,
					)
					return nil
				}

				if len(proxies) == 0 {
					disconnectBecauseNotHaveProxyToSend(
						connection,
					)
					return nil
				}

				connection.PacketQueue[1] = loginStart
				connection.Stop = true

				key, err := ProxyApp.GetRandomProxy(proxies)

				if err != nil {
					disconnectBecauseNotHaveProxyToSend(
						connection,
					)
					return nil
				}

				go ProxyConnector.ConnectToProxy(connection, key)
			} else {
				fmt.Printf("Falha ao receber a conexão de %s", string(loginStart.Username))
			}
		}
	default:
		{
			fmt.Println("Não foi possível ler esse estado", connection.State)
		}
	}

	return nil
}

func disconnectBecauseNotHaveProxyToSend(connection *protocol.Connection) {
	connection.Disconnect(chat.TextComponent{
		Text: fmt.Sprintf(
			"%s\n\n§r§cNão foi possível localizar um proxy para enviar você.",
			Constants.SERVER_PREFIX,
		),
	})
}

func disconnectBecauseMaintenanceModeIsEnabled(connection *protocol.Connection) {
	connection.Disconnect(chat.TextComponent{
		Text: fmt.Sprintf(
			"%s\n\n§r§cO servidor atualmente encontra-se em manutenção.",
			Constants.SERVER_PREFIX,
		),
	})
}