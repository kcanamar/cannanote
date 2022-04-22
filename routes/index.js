////////////////////////
// Setup - Import deps
////////////////////////
const express = require('express')
const router = express.Router()

///////////////////////
// Declare Routes 
///////////////////////
router.get('/', function(req, res, next) {
    res.redirect('/cannanote')
})
///////////////////////
// Exports
///////////////////////
module.exports = router