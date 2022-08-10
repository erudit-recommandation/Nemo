clear-docker:
	docker system prune -a -f

full-clear-docker:
	docker system prune -a -f --volumes

run-docker:
	sudo chmod -R a+rwx ./data/
	docker-compose up --build

run-docker-debug:
	sudo chmod -R a+rwx ./data/
	docker-compose -f "docker-compose-debug.yml" up --build 

see-log-web-app:
	docker-compose logs erudit_recommandation

stop-docker:
	docker-compose stop

run-text-analysis-service:
	cd ./text_analysis_service && export FLASK_APP=app && flask run

delete-database:
	sudo  rm -r ./data/arango/
	
run-etl:
	cd initialisation-service && make run

test:
	$(GOCMD) test $(TEST_PATH)
