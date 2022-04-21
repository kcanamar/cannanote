////////////////////////
// Setup - Import deps
////////////////////////
const express = require('express')
const router = express.Router()
///////////////////////
// Declare Routes and Routers 
///////////////////////
// INDUCES - Index, New, Delete, Update, Create, Edit, Show
// Index
router.get("/", (req, res) => {
    res.send('It works')
})
// New
router.get("/new", (req, res) => {
    res.send('Something New')
})
// Edit
router.get("/:id/edit", (req, res) => {
    res.send(`here is what you asked to edit ${req.params.id}`)
})
// Show
router.get("/:id", (req, res) => {
    res.send(`Order for ${req.params.id}`)
})
// Update
router.put("/:id", (req, res) => {
    console.log(`here is what you wanted ${req.params.id}`)
    res.send(`It is Changed ${req.body}`)
})
// Create 
router.post("/", (req, res) => {
    res.send(`here is what you made ${req.body}`)
})
// Delete
router.delete("/:id", (req, res) => {
    res.send(`${req.params.id}....Target neutralized`)
})
///////////////////////
// Exports
///////////////////////
module.exports = router