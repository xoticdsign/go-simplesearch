package main

import (
	"github.com/xoticdsign/go-simplesearch/internal/app"
)

/*

THIS PROJECT WAS CREATED AS AN EXAMPLE IMPLEMENTATION OF ELASTICSEARCH IN GOLANG.

Overview:
-----------
The purpose of this project is to demonstrate a simple implementation of Elasticsearch using Golang. It showcases how to integrate a search engine (in this case, Elasticsearch) with a Go web application, handling requests, performing search queries, and returning results. The application utilizes an API endpoint for search requests and is structured in a way that allows easy configuration, error handling, and graceful shutdown.

Main Features:
--------------
1. Search Functionality:
   - The core functionality revolves around searching for products based on a request from the client. A user can provide search parameters (like search terms and filters) and receive relevant product results.

2. Elasticsearch Integration:
   - Elasticsearch is used as the search engine for querying product data. The application connects to Elasticsearch and performs searches based on the input from the client. The search results are returned to the user as JSON responses.

3. Graceful Shutdown:
   - The application handles signals like SIGINT (Ctrl+C) or SIGTERM gracefully, ensuring that ongoing processes are cleaned up before the app shuts down, providing a smoother user experience during restarts or shutdowns.

4. Configurable Environment:
   - The application reads configuration settings such as host, port, timeouts, and Elasticsearch credentials from a YAML configuration file or environment variables. This makes the app flexible and adaptable to different deployment environments like local development, testing, and production.

How it Works:
-------------
1. Environment Setup:
   - The `env` environment variable determines the configuration file to be loaded. If not specified, the app defaults to "local" settings, which typically point to a development environment.

2. App Initialization:
   - The `app.New` function initializes the application by loading the configuration and setting up necessary services, such as the Elasticsearch connection and the web server.

3. Running the App:
   - Once the app is initialized, it is started by calling the `Run()` method. This method will start the web server and wait for incoming requests. It also listens for shutdown signals (SIGINT or SIGTERM) to gracefully terminate the application.

4. Search Endpoint:
   - The main feature is a `/search` endpoint that listens for POST requests. Clients can send search queries, which the application then processes and queries Elasticsearch to return search results.

5. Error Handling:
   - If any error occurs during the app initialization, running, or processing, the application panics and prints the error. This is a simple way of handling critical failures, but more sophisticated error handling and logging could be added for production environments.

In summary, this project provides a straightforward example of how to build a Go web application with Elasticsearch as the backend for search functionality, while emphasizing flexibility, error handling, and graceful shutdown mechanisms.

*/

func main() {
	a, err := app.New()
	if err != nil {
		panic(err)
	}
	err = a.Run()
	if err != nil {
		panic(err)
	}
}
