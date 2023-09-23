package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	_ = godotenv.Load()
)

// //////////////////////////////////////////////////////////////////////////
// usage:
//
//	go run main.go gql --in <gql filename>
func main() {
	ctx := context.Background()

	gqlCmd := flag.NewFlagSet("gql", flag.ExitOnError)
	gqlCmdFilenameInPtr := gqlCmd.String("in", "", "GraphGQ File")

	syncKcUsersCmd := flag.NewFlagSet("syncusers", flag.ExitOnError)
	syncKcUsersBackupDbDir := syncKcUsersCmd.String("dbdir", "", "Local DB dir for backup data")
	if len(os.Args) < 2 {
		fmt.Println("expected 'gql' or 'synckcusers' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "syncusers":
		syncKcUsersCmd.Parse(os.Args[2:])
		if len(*syncKcUsersBackupDbDir) == 0 {
			*syncKcUsersBackupDbDir = os.Getenv("FRDB_USERS_BACKUP_DIR")
		}
		if len(*syncKcUsersBackupDbDir) == 0 {
			log.Panic("--dbdir or env:FRDB_USERS_BACKUP_DIR is required")
		}
		SyncKcUsers(ctx, syncKcUsersBackupDbDir)
	case "gql":
		gqlCmd.Parse(os.Args[2:])
		if 0 >= len(*gqlCmdFilenameInPtr) {
			log.Panic("in param required for gql request")
		}
		MakeGqlReq(ctx, gqlCmdFilenameInPtr)
	case "gentoken":
		_, token := LoginKcAdmin(ctx)
		log.Printf("Bearer %s", token)
	default:
		flag.PrintDefaults()
		log.Fatal("Invalid cli flag combination")
		os.Exit(1)
	}

	log.Println("Done")
}
