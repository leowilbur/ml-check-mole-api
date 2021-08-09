package auth

var permissionsMap = map[string][]string{
	"Default": []string{
		"owned.lesions.read",
		"owned.lesions.write",
		"owned.requests.create",
	},
	"Doctors": []string{
		"requests.read",
		"reports.read",
		"lesions.read",
		"requests.respond",
	},
	"Administrators": []string{
		"body-parts.write",
		"questions.write",
		"lesions.read",
		"lesions.write",
		"requests.create",
		"requests.respond",
	},
}

// LookupPermissions matches user's groups and checks if the user has permissions
// passed to the function.
func LookupPermissions(groups []string, permissions ...string) bool {
	lookup := map[string]struct{}{}
	for _, item := range permissionsMap["Default"] {
		lookup[item] = struct{}{}
	}

	for _, group := range groups {
		slice, ok := permissionsMap[group]
		if !ok {
			continue
		}

		for _, item := range slice {
			lookup[item] = struct{}{}
		}
	}

	for _, permission := range permissions {
		if _, ok := lookup[permission]; !ok {
			return false
		}
	}

	return true
}
