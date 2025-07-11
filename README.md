# Implementation of Test Task

## Description

### Test task: Web-application for analyzing Webpages
**Objective:**
The objective is to build a web application that does an analysis of a web-page/URL.
The application should show a form with a text field in which users can type in the URL of the webpage to be analyzed. Additionally to the form, it should contain a button to send a request to the server.
After processing the results should be shown to the user.

#### Results should contain next information:
* What HTML version has the document?
* What is the page title?
* How many headings of what level are in the document?
* How many internal and external links are in the document?Are there any inaccessible links and how many?
* Does the page contain a login form?
* In case the URL given by the user is not reachable an error message should be presented to a user. The message should contain the HTTP status code and a useful error description.

## How to Run Locally

A Makefile is provided to simplify launching the application on your local machine. You need to have GNU `make` installed. There are two ways to run the application:

* **In a Docker container** — Docker must be installed.
* **Directly on your localhost** — Golang must be installed.

### Run the application in Docker

1. Copy `.env.dist` to `.env`
2. Run:
   ```bash
   make docker-build
   make docker-run
   ```

### Run the HTTP server with Go
1. Copy .env.dist into .env
2. Run
```bash
  make go-run
```
or
```bash
  go run main.go start
```

### Build application with Golang
1. Copy .env.dist into .env
2. To build `app-test` in the project directory run
```bash
  make go-build
```
or
```bash
  go build -v -o app-test
```

### Run HTTP server with built application
```bash
  ./app-test start
```

## Configuration
Configuration of the server host and port could be changed in `.env` file


## Assumptions
### Point from the description: **"Are there any inaccessible links and how many?"**
A link is considered inaccessible if the target page:
1. Does not respond within a predefined timeout 10 seconds
2. Returns an HTTP status code 400 and more.

### Point from the description: **"Does the page contain a login form?"**
Login form is considered to be present if the page contains a form with an input field with `type="password"`.

## It's worth to implement
1. BDD scenarios for HTPP server are absent. It's need to add them to check application behavior.
2. Javascript is not used there. Every submit reboot th web page. It would be better to do async calls from FE and receive data from BE in JSON.
3. Timeouts for HTTP requests are hardcoded in the code. It's need to add them to the configuration file. It will decrease time of testing.
4. HTTP client doesn't send headers as a browser and bot-protected URL is not available. It worth to add them to overcome the bot-protection.
