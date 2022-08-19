package psql

const (
	AddUser = `insert into public.users (id, fullname, email, username, password, token, refreshtoken, created, updated)
			values ('%s', '%s', '%s','%s', '%s', '%s', '%s', '%s', '%s');`

	FindUserByUsername = `select id, fullname, email, username, password, token, created, updated, refreshtoken
			from public.users
			where username = '%s';`

	FindUserById = `select id, fullname, email, username, password, token, created, updated, refreshtoken
			from public.users
			where id = '%s';`

	UpdateTokens = `update public.users
			set token = '%s', refreshtoken = '%s', updated = '%s'
			where id = '%s';`
)
