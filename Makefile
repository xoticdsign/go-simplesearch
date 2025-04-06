# VARIABLES #######################################################################################################################

simplesearch := cmd/simplesearch/main.go 															       # SimpleSearch main()
esmigrator := cmd/esmigrator/main.go																	   # ESMigrator main()

esaddress := https://0.0.0.0:9200																	       # ElasticSearch Address
esusername := xoticdsign                                                                                   # ElasticSearch Username
espassword := 188696                                    											       # ElasticSearch Password

esmigrations := ./migrations/elasticsearch/migration.json 									     	       # Migrations Location

# SIMPLE SEARCH APP ###############################################################################################################

run_simplesearch: $(simplesearch)
	ES_USERNAME=$(esusername) ES_PASSWORD=$(espassword) go run $(simplesearch)

# TOOLS ###########################################################################################################################

esmigrate_up: $(esmigrator)
	DIRECTION=up MIGRATIONS=$(migrations) ADDRESS=$(esaddress) USERNAME=$(esusername) PASSWORD=$(espassword) go run $(esmigrator)

esmigrate_down: $(esmigrator)
	DIRECTION=down MIGRATIONS=$(migrations) ADDRESS=$(esaddress) USERNAME=$(esusername) PASSWORD=$(espassword) go run $(esmigrator)