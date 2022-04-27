////////////////////////
// Setup - Import deps
////////////////////////
const express = require('express')
const router = express.Router()
const entryCtrl = require('../controllers/entries')
///////////////////////
// Declare Routes and Routers 
///////////////////////
// INDUCES - Index, New, Delete, Update, Create, Edit, Show
router.use((req, res, next) => {
    console.log("canna router middle ware" + req.session.loggedIn)
    if (req.session.loggedIn) {
        next()
    } else {
        res.redirect("/")
    }
})
router.get("/", entryCtrl.index)
router.get("/new", entryCtrl.new)
router.get("/seed", entryCtrl.seed)
router.get("/:id/edit", entryCtrl.edit)
router.get("/:id", entryCtrl.show)
router.put("/:id/like", entryCtrl.like)
router.put("/:id", entryCtrl.update)
router.post("/", entryCtrl.create)
router.delete("/:id", entryCtrl.delete)
///////////////////////
// Exports
///////////////////////
module.exports = router