
# Product Requirements Documentation

**Summary**
| Field | Detail |
|-------|--------|
| Project Name | CannaNote |
| Description | Journal for Cannabis Users |
| Developers | Kyle Canamar |
| Live Website | [CannaNote](https://karc-cannanote.herokuapp.com) |
| Repo | [GitHub](https://github.com/kcanamar/cannanote) |

## Problem Being Solved and Target Market

CannaNote is a social site where you can share your experience with Cannabis. Everyone is different and cannabis affects us each in different ways, this service will allow users to share thier unique experiences.

## On the Horizion

- Age Verification
- Inhanced user experience
- Seach functionality
- Sort functionality
- Noticifcations

## User Stories

List of stories users should experience when using your application.

- Users should be able to see the site on desktop and mobile
- Users can create an account
- Users can sign in to their account
- Users can create a new Entry
- Users can see all Entries in the Feed
- Users can update Entries
- User can delete Entries

## Route Tables

- The endpoint: the URL to which the request must be made
- The method: the type of http method the request should be
- The response: what the response should be, a web page, json data, etc.

| Endpoint | Method | Response | Other |
| -------- | ------ | -------- | ----- |
| /Cannanote | GET | Brings up the Feed  | |
| /Cannanote | POST | Create new Entry Posting to the top of Feed | body must include data for new item |
| /Cannanote/:id/edit | GET | Shows the Entry edit page | |
| /Cannanote/:id | GET | Shows the Entry in greater detail | |
| /Cannanote/:id | PUT | update item with matching idea, return to the Feed | body must include updated data |
| /Cannanote/:id/like | PUT | Adds one like, return to the Feed | |
| /Cannanote/:id | DELETE | delete the Enrty with the matching mongoDB id | |
| / | GET | Displays the Entrance to app | |
| /signup | POST | creates new user account returns back to login screen | new user info must be included in body |
| /login | POST | logs in user and returns user with session cookie | username and password must be included in body |

## User Interface Mockups
|:-------------------------|:-------------------------|:--------Wire Frames--------|:-------------------------|:-------------------------|
|Login/Signup Pages <img src="./public/wireframes/Meen Project Login-signup.png" alt=""></img> |Home Page <img src="./public/wireframes/Meen Project Home page.png" alt=""></img> |Create/Update Pages <img src="./public/wireframes/Meen Project Create-Update.png" alt=""></img> | | |
|:-------------------------|:-------------------------|:--------Model Schema's--------|:-------------------------|:-------------------------|
|User And Entries <img src="./public/wireframes/Model Schema's.png" alt=""></img> | | | | |
|:-------------------------|:-------------------------|:--------Screen Shots--------|:-------------------------|:-------------------------|
|Main <img src="./public/screenshots/main.png" alt=""></img> | New User<img src="./public/screenshots/new-user.png" alt=""></img> | Existing User <img src="./public/screenshots/existing-user.png" alt=""></img> | Feed <img src="./public/screenshots/feed.png" alt=""></img> | New Entry <img src="./public/screenshots/new.png" alt=""></img> |
| Create <img src="./public/screenshots/create.png" alt=""></img> | Show <img src="./public/screenshots/show.png" alt=""></img> | Before Delete <img src="./public/screenshots/before-delete.png" alt=""></img> | After Delete <img src="./public/screenshots/after-delete.png" alt=""></img> | Mobile Test <img src="./public/screenshots/mobile-test.png" alt=""></img> |
