package users

import (
	"fmt"
	cache2 "github.com/patrickmn/go-cache"
	"log"
	Databases "net/hyren/nyrah/databases"
	"time"
)

var (
	CACHE = cache2.New(cache2.NoExpiration, 10*time.Second)
)

func IsHelperOrHigher(name string) bool {
	userGroupsDue, found := CACHE.Get(fmt.Sprintf("is_helper_%s", name))

	if !found {
		db := Databases.StartMariaDB()

		rows, err := db.Query(
			fmt.Sprintf(
				"SELECT User.`id`, UserGroupDue.`group_name`, UserGroupDue.`due_at` FROM `users` AS User INNER JOIN `users_groups_due` AS UserGroupDue WHERE User.name LIKE '%s' AND UserGroupDue.user_id=User.id;",
				name,
			),
		)

		defer db.Close()

		if err != nil {
			log.Println(err)
		}

		userGroupsDue = rows.Next()

		defer rows.Close()

		CACHE.Set(fmt.Sprintf("is_helper_%s", name), userGroupsDue, 5*time.Minute)
	}

	return userGroupsDue.(bool)
}
