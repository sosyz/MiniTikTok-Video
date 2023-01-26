package conf

type Redis struct {
	Host     string
	Port     int
	Password string
	Db       int
}

type Pg struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
	Db       string
	TimeZone string
}

type Neo4j struct {
	Host     string
	Port     int
	User     string
	Password string
	Realm    string
}

type Listen struct {
	Host string
	Port int
}

type Node struct {
	Id int
}

type Server Listen

type Run struct {
	Listen Listen
	Server Server
}
