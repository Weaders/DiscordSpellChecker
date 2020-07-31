package speller

var usersCount map[string]int = make(map[string]int)

// Add val for user counter
func AddCountForUser(userID string, count int) {

	val, ok := usersCount[userID]

	if ok {
		usersCount[userID] = val + count
	} else {
		usersCount[userID] = count
	}

}

//Return value of counter for specific users
func CounterForUsers(usersIDs []string) map[string]int {

	result := make(map[string]int)

	for _, v := range usersIDs {

		val, ok := usersCount[v]

		if ok {
			result[v] = val
		}

	}

	return result

}
