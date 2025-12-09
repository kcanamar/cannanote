
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

CannaNote is a social site where you can share your experience with Cannabis. Everyone is different and cannabis affects us each in different ways, this service will allow users to share their unique experiences.

## On the Horizon

- Age Verification
- Medical Cannabis Patient Focus
- Search & Sort functionality
- Notifications

## Getting Started

- TODO: update on boarding documentation

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
| /Cannanote/:id | DELETE | delete the Entry with the matching mongoDB id | |
| / | GET | Displays the Entrance to app | |
| /signup | POST | creates new user account returns back to login screen | new user info must be included in body |
| /login | POST | logs in user and returns user with session cookie | username and password must be included in body |
