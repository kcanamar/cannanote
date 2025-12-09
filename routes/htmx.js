////////////////////////
// Setup - Import deps
////////////////////////
const express = require('express')
const router = express.Router()
const htmxCtrl = require('../controllers/htmx')
///////////////////////
// Declare Routes and Routers 
///////////////////////
// INDUCES - Index, New, Delete, Update, Create, Edit, Show
router.use((req, res, next) => {
    if (req.session.loggedIn) {
        next()
    } else {
        res.redirect("/")
    }
})
router.get("/deleteConfirm", htmxCtrl.deleteConfirm)
///////////////////////
// Exports
///////////////////////
module.exports = router
