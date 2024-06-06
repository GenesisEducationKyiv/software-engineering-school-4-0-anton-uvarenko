# backend_school
In order to run this test task you need to edit .env.prod file.  
But before that go to your gmail account and turn on 2 factor auth as said here https://support.google.com/accounts/answer/185833?hl=en
`FROM_EMAIL` is value that represents an email from which you gonna send messages to subscribed emails.  
`FROM_EMAIL_PASSWORD` is value that you've generated in 2 factor auth in your google profile.
    
In this project i tried to implement clean architecture. It is incomplete i don't have repositories, cause i use sqlc. When u run `docker compose up -d` i have two steps for building application.  
First builds all executables  
Second copies builds with migrations, than runs migrations and starts the server.  

Before starting the server application builds all the required parts.  
Firstly is opens db connection in db package  
Then in creates instances of currenyConverter(uses monobank api) and emailSender(uses gomail).  
Then builds service and handler.  
Then builds server and sets up routes.  
Then builds scheduler for cronjobs.
Then starts the app.

