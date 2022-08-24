package filerecord

type FileRecord struct {
	uuid            string
	hash            string
	filename        string
	extension       string
	version         int
	literalLocation string
	logicalLocation string
	dateEntered     string
	dateCreated     string
	dateModified    string
}
