package psql

const (
	AddUser = `insert into public.users (id, fullname, email, username, password, token, refreshtoken, created, updated, shows)
			values ('%s', '%s', '%s','%s', '%s', '%s', '%s', '%s', '%s', '{%s}') returning id;`

	FindUserByEmail = `select id, fullname, email, username, password, token, refreshtoken, created, updated, shows
			from public.users
			where email = '%s';`

	FindUserById = `select id, fullname, email, username, password, token, refreshtoken, created, updated, shows
			from public.users
			where id = '%s';`

	UpdateTokens = `update public.users
			set token = '%s', refreshtoken = '%s', updated = '%s'
			where id = '%s';`

	SetUserShows = `update public.users set shows='{%s}'
			where id = '%s';`

	AddUserShow = `update public.users 
			set shows = array_append(shows, '%s')
			where id = '%s';`
	//AddUserShow = `select array(select distinct unnest(shows || '{%s}'))
	//		from public.users
	//		where id = '%s';`
)
