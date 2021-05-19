package users

import (
	"fmt"
	"log"
	"net/hyren/nyrah/cache/local"
	"net/hyren/nyrah/minecraft/chat"
	"net/hyren/nyrah/minecraft/protocol"
	"net/hyren/nyrah/misc/providers"
	"time"

	Constants "net/hyren/nyrah/misc/constants"
)

func IsHelperOrHigher(name string) bool {
	userGroupsDue, found := local.CACHE.Get(fmt.Sprintf("is_helper_or_higher_%s", name))

	if !found {
		connection := providers.POSTGRESQL_MAIN.Provide()

		defer connection.Close()

		rows, err := connection.Query(fmt.Sprintf(
				`SELECT "users"."id", "users_groups_due"."group_name", "due_at" FROM "users" INNER JOIN "users_groups_due" ON "user_id"="users"."id" AND "users"."name" ILIKE '%s' AND "users_groups_due"."group_name"=ANY(ARRAY['MASTER', 'DIRECTOR', 'MANAGER', 'MODERATOR', 'HELPER']);`,
				name,
		))

		if err != nil {
			log.Println(err)
		}

		userGroupsDue = rows.Next()

		defer rows.Close()

		local.CACHE.Set(fmt.Sprintf("is_helper_or_higher_%s", name), userGroupsDue, 5*time.Minute)
	}

	return userGroupsDue.(bool)
}

func DisconnectBecauseNotHaveProxyToSend(connection *protocol.Connection) {
	connection.Disconnect(chat.TextComponent{
		Text: fmt.Sprintf(
			"%s\n\n§r§cNão foi possível localizar um proxy para enviar você.",
			Constants.SERVER_PREFIX,
		),
	})
}

func DisconnectBecauseMaintenanceModeIsEnabled(connection *protocol.Connection) {
	connection.Disconnect(chat.TextComponent{
		Text: fmt.Sprintf(
			"%s\n\n§r§cO servidor atualmente encontra-se em manutenção.",
			Constants.SERVER_PREFIX,
		),
	})
}