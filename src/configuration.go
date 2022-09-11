// Configuration class that is used to store information for the application. 
// This is expected to be created before pretty much anything else when the app starts, 
// and should be used to initialize connections out (database), and service config 

package main 


// TODO: annotate fields and sort out loading from a config file. 

type LibraryConfig struct {
    ArchiveBasePath string 
    DbConnectionString string 
}

func NewLibraryConfig(archiveBasePath string, dbConnectionString string) *LibraryConfig {
    return &LibraryConfig{
        ArchiveBasePath: archiveBasePath,
        DbConnectionString: dbConnectionString,
    }
}
