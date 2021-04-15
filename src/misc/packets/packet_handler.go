package packets

import (
	"errors"
	"fmt"
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
				loginSuccess, _ := holder.(packet.LoginSuccess)

				log.Println(loginSuccess.UUID)

				name := string(loginStart.Username)

				log.Println("Conexão recebida de", name)

				if Config.IsMaintenanceModeEnabled() == true && !canJoin(
					name,
				) {
					disconnectBecauseMaintenanceModeIsEnabled(
						connection,
					)
					return nil
				}

				db := Databases.StartPostgres()

				rows, err := db.Query("SELECT \"name\" FROM \"applications\" WHERE \"application_type\"='PROXY';")

				defer db.Close()

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
	db := Databases.StartPostgres()

	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT \"id\" FROM \"users\" WHERE \"name\" ILIKE '%s';",
			name,
		),
	)

	defer db.Close()

	if err != nil {
		log.Println(err)

		defer rows.Close()
		return false
	}

	var id string

	if rows.Next() {
		rows.Scan(&id)
	}

	defer rows.Close()

	db = Databases.StartPostgres()

	rows, err = db.Query(
		fmt.Sprintf(
			"SELECT \"group_name\" FROM \"users_groups_due\" WHERE \"user_id\"='%s';",
			id,
		),
	)

	defer db.Close()

	if err != nil {
		log.Println(err)

		defer rows.Close()
		return false
	}

	for next := rows.Next(); next; next = rows.Next() {
		var group_name string

		rows.Scan(&group_name)

		if Config.IsGroupWhitelisted(group_name) {
			return true
		}
	}

	defer rows.Close()
	defer db.Close()

	return false
}
