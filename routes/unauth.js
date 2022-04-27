////////////////////////
// Setup - Import deps and create app object
////////////////////////
const express = require('express')
const router = express.Router()
const unauthCtrl = require('../controllers/unauth')
///////////////////////
// Declare Routes and Routers 
///////////////////////
// INDUCES - Index, New, Delete, Update, Create, Edit, Show
router.get("/", unauthCtrl.main)
router.post("/signup", unauthCtrl.signup)
router.post("/login", unauthCtrl.login)
router.post("/logout", unauthCtrl.logout)
///////////////////////////
// Server Listener
/////////////////////////
module.exports = router