////////////////////////
// Setup - Import deps
////////////////////////
const express = require('express')
const methodOverride = require('method-override')
const morgan = require('morgan')

///////////////////////////
// Export
///////////////////////////
module.exports = function(app) {
    app.use(methodOverride('_method'))
    app.use("/static", express.static('public'))
    app.use(express.urlencoded({extended: true}))
    app.use(morgan('tiny'))
}