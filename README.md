# Chat System Go APIs

### Overview
This is a creation of messages and chats endpoints in goLang to achieve high performance and concurrency.
### Goals
* Create chat and message Endpoints.
* Send chat or message create payload to rabbitMQ.
* Handle concurrent application runs and race condition.
* Containerize the task.
### Requirements
In order to run the application please follow the steps:
- Ubuntu running OS.
- Docker installed.
### Quick start
* Clone the project.
* run docker-compose up the docker file.

### Test services
* ##### Postman or curl command line
  *`/api/v1/applications/{token}/chats [POST] {"chat_name": "chat_name"}`\
  *`/api/v1/applications/{token}/chats/{chat_id}/messages [POST] {"message":"message"}`
  



