# COMMON #########################################################################################################################################################################################

es_address := https://0.0.0.0:9200					
es_username := xoticdsign                              
es_password := 188696                                

# SIMPLE SEARCH APP ##############################################################################################################################################################################

simplesearch := cmd/simplesearch/main.go 			  

simplesearch_address := 0.0.0.0:8080						

simplesearch: $(simplesearch)
	ADDRESS=$(address) ES_ADDRESS=$(es_address) ES_USERNAME=$(es_username) ES_PASSWORD=$(es_password) go run $(simplesearch)

# TOOLS ##########################################################################################################################################################################################

esmigrator := cmd/esmigrator/main.go					    

esmigrator_direction := up									
esmigrator_migrations := ./migrations/elasticsearch/migration.json 	
esmigrator_migrate_down_with_index := false						 

esmigrator: $(esmigrator)
	DIRECTION=$(esmigrator_direction) MIGRATIONS=$(esmigrator_migrations) MIGRATE_DOWN_WITH_INDEX=$(esmigrator_migrate_down_with_index) ES_ADDRESS=$(es_address) ES_USERNAME=$(es_username) ES_PASSWORD=$(es_password) go run $(esmigrator)

# DOCKER #########################################################################################################################################################################################

docker_container_name := simplesearch
docker_container_port := 8080:8080
docker_image_name := simplesearch
docker_dockerfile := deployments/docker/Dockerfile
docker_env_es_address := https://host.docker.internal:9200

docker_build:
	docker build -f $(docker_dockerfile) -t $(docker_image_name) .   

docker_run:
	docker run --name $(docker_container_name) -p $(docker_container_port) -e ADDRESS=$(simplesearch_address) -e ES_ADDRESS=$(docker_env_es_address) -e ES_USERNAME=$(es_username) -e ES_PASSWORD=$(es_password) $(docker_image_name)