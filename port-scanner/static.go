package portscanner

var HighSeverityPorts []int = []int{
	15,   // Netstat
	20,   // FTP
	21,   // FTP
	22,   // SSH
	23,   // Telnet Protocol
	25,   // SMTP
	53,   // DNS
	69,   // TFTP
	88,   // Kerberos
	110,  // POP3
	111,  // Port Map
	143,  // IMAP
	161,  // SNMP
	389,  // LDAP
	1521, // Oracle Database
	1830,
	2424, // OrientDB
	2501, // Standard Port Monitor (using RAW) - Printer
	3306, // MariaDB, MySQL
	3389, // RDP
	5001, // Intelligent Printer Data Stream - Printer
	5432, // PostgreSQL
	5984, // CouchDB
	6001, // Printer
	6379, // Redis
	7210,
	7687,  // Neo4j
	8086,  // InfluxDB
	8529,  // ArangoDB
	9001,  // Standard Port Monitor (using RAW) - Printer
	9002,  // Standard Port Monitor (using RAW) - Printer
	9003,  // Standard Port Monitor (using RAW) - Printer
	9004,  // Standard Port Monitor (using RAW) - Printer
	9005,  // Standard Port Monitor (using RAW) - Printer
	9006,  // Standard Port Monitor (using RAW) - Printer
	9007,  // Standard Port Monitor (using RAW) - Printer
	9600,  // Intelligent Printer Data Stream - Printer
	26257, // CockroachDB
	27017, // MongoDB
	27018, // mongod
	27019, // mongod
	28015, // RethinkDB
	28017, // MongoDb Http interface
}
