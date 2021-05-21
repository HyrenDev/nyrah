package packets

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/hyren/nyrah/minecraft/protocol"
	"net/hyren/nyrah/minecraft/protocol/codecs"
	"net/hyren/nyrah/minecraft/protocol/packet"
	"reflect"
	"strings"

	ProxyApplication "net/hyren/nyrah/applications"
	ProxyStatus "net/hyren/nyrah/applications/status/proxy"
	ProxyConnector "net/hyren/nyrah/misc/connector"
	NyrahConstants "net/hyren/nyrah/misc/constants"
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
			handshake.ServerAddress = codecs.String(strings.Split(connection.Handle.RemoteAddr().String(), ":")[0])

			log.Println("Received handshake request from:", connection.Handle.RemoteAddr().String())

			connection.PacketQueue[0] = handshake

			return nil
		}
	case protocol.Status:
		{
			_, ok := holder.(packet.StatusRequest)

			var online = ProxyStatus.GetOnlinePlayers()
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
				name := string(loginStart.Username)

				log.Println(fmt.Sprintf(
					"Conexão recebida de [%s/%s]",
					name,
					connection.GetRemoteAddr(),
				))

				if Config.IsMaintenanceModeEnabled() == true && !User.IsHelperOrHigher(
					name,
				) {
					User.DisconnectBecauseMaintenanceModeIsEnabled(
						connection,
					)
					return nil
				}

				proxies, err := ProxyApplication.FetchAvailableProxiesNames()

				if err != nil {
					User.DisconnectBecauseNotHaveProxyToSend(
						connection,
					)
					return nil
				}

				if len(proxies) == 0 {
					User.DisconnectBecauseNotHaveProxyToSend(
						connection,
					)
					return nil
				}

				connection.PacketQueue[1] = loginStart
			}

			encryptionRequest, ok := holder.(packet.EncryptionRequest)

			if ok {
				privateKey, err := rsa.GenerateKey(rand.Reader, 1024)

				if err != nil {
					panic(err)
				}

				publicKey := privateKey.PublicKey

				publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)

				token := make([]byte, 4)

				verifyTokenLength, err := rand.Read(token)

				if err != nil {
					panic(err)
				}

				encryptionRequest.ServerID = NyrahConstants.SERVER_ID
				encryptionRequest.PublicKey = publicKeyBytes
				encryptionRequest.PublicKeyLength = codecs.VarInt(len(publicKeyBytes))
				encryptionRequest.VerifyToken = token
				encryptionRequest.VerifyTokenLength = codecs.VarInt(verifyTokenLength)

				connection.PacketQueue[2] = encryptionRequest
			}

			encryptionResponse, ok := holder.(packet.EncryptionResponse)

			if ok {
				response, err := http.Get(fmt.Sprintf(
					`https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s&ip=%s`,
					string(loginStart.Username),
					encryptionRequest.ServerID,
					connection.GetRemoteAddr(),
				))

				if err != nil {
					return nil
				}

				encryptionResponse.SharedSecret = encryptionRequest.PublicKey
				encryptionResponse.SharedSecretLength = encryptionRequest.PublicKeyLength
				encryptionResponse.VerifyToken = encryptionRequest.VerifyToken
				encryptionResponse.VerifyTokenLength = encryptionRequest.VerifyTokenLength

				connection.PacketQueue[5] = encryptionResponse

				println(response.Body)

				response.Close = true
			}

			loginSuccess, ok := holder.(packet.LoginSuccess)

			if ok {
				connection.PacketQueue[7] = loginSuccess
				connection.Stop = true

				key, err := ProxyApplication.GetRandomProxy()

				if err != nil {
					User.DisconnectBecauseNotHaveProxyToSend(
						connection,
					)
					return nil
				}

				go ProxyConnector.ConnectToProxy(connection, key)
			}
		}
	default:
		{
			log.Println("Não foi possível ler esse estado", connection.State)
		}
	}

	return nil
}
