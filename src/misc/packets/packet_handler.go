package packets

import (
	"encoding/hex"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gominet/chat"
	"gominet/protocol"
	"gominet/protocol/codecs"
	"gominet/protocol/packet"
	"log"
	"reflect"
	"strings"

	ProxyApp "../../applications"
	Databases "../../databases"
	Connection "../connection"
	NyrahConstants "../constants"
	Config "../utils"
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
			handshake.ServerAddress = codecs.String(
				string(
					handshake.ServerAddress,
				) + "%ABC%" + strings.Split(connection.Handle.RemoteAddr().String(), ":")[0],
			)

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
					log.Println("Cannot find favicon")
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
				if Config.IsMaintenanceModeEnabled() == true && !canJoin(
					string(loginStart.Username),
				) {
					disconnectBecauseMaintenanceModeIsEnabled(
						connection,
					)
					return nil
				}

				db := Databases.StartPostgres()

				rows, err := db.Query("SELECT \"name\" FROM \"applications\" WHERE \"application_type\"='PROXY';")

				if err != nil {
					return err
				}

				var proxies []string

				var index = 0

				for rows.Next() {
					var name string

					err := rows.Scan(&name)

					if err != nil {
						return err
					}

					proxies = append(proxies, name)

					index++
				}

				defer rows.Close()
				defer db.Close()

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

				go Connection.SendToProxy(connection, key)
			}
		}
	default:
		{
			//
		}
	}

	return nil
}

func disconnectBecauseNotHaveProxyToSend(connection *protocol.Connection) {
	connection.Disconnect(chat.TextComponent{
		Text: fmt.Sprintf(
			"%s\n\n§r§cNão foi possível localizar um proxy para enviar você.",
			NyrahConstants.SERVER_PREFIX,
		),
	})
}

func disconnectBecauseMaintenanceModeIsEnabled(connection *protocol.Connection) {
	connection.Disconnect(chat.TextComponent{
		Text: fmt.Sprintf(
			"%s\n\n§r§cO servidor atualmente encontra-se em manutenção.",
			NyrahConstants.SERVER_PREFIX,
		),
	})
}

func canJoin(name string) bool {
	userId, err := offlinePlayerUUID(name)

	if err != nil {
		return false
	}

	db := Databases.StartPostgres()

	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT \"group_name\" FROM \"users_groups_due\" WHERE \"user_id\"='%s';",
			userId,
		),
	)

	if err != nil {
		return false
	}

	for next := rows.Next(); next; next = rows.Next() {
		var group_name string

		_ = rows.Scan(&group_name)

		log.Println("Grupo de ", name, " -> ", group_name)

		if group_name == "MASTER" || group_name == "DIRECTOR" || group_name == "MANAGER" || group_name == "ADMINISTRATOR" || group_name == "MODERATOR" || group_name == "HELPER" {
			return true
		}
	}

	return false
}

func offlinePlayerUUID(name string) (uuid.UUID, error) {
	if len(name) == len(uuid.Nil.String()) {
		return uuid.FromString(name)
	}

	b, err := hex.DecodeString(name)

	if err != nil {
		return uuid.Nil, err
	}

	return uuid.FromBytes(b)
}
