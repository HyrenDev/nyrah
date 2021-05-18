package users

import (
	"fmt"
	"net/hyren/nyrah/cache/local"
	"time"

	NyrahProvider "net/hyren/nyrah/misc/providers"
)

func IsHelperOrHigher(name string) bool {
	userGroupsDue, found := local.CACHE.Get(fmt.Sprintf("is_helper_or_higher_%s", name))

	if !found {
		rows, err := NyrahProvider.MARIA_DB_MAIN.Provide().Query(
			fmt.Sprintf(
				"SELECT User.`id`, UserGroupDue.`group_name`, UserGroupDue.`due_at` FROM `users` AS User INNER JOIN `users_groups_due` AS UserGroupDue WHERE User.name LIKE '%s' AND UserGroupDue.user_id=User.id;",
				name,
			),
		)

		if err != nil {
			fmt.Println(err)
		}

		userGroupsDue = rows.Next()

		defer rows.Close()

		local.CACHE.Set(fmt.Sprintf("is_helper_or_higher_%s", name), userGroupsDue, 5*time.Minute)
	}

	return userGroupsDue.(bool)
}
