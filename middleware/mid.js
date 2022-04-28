////////////////////////
// Setup - Import deps
////////////////////////
require("dotenv").config()
const express = require('express')
const methodOverride = require('method-override')
const morgan = require('morgan')
const session = require("express-session")
const MongoStore = require("connect-mongo")
const Entry = require("../models/entries")
const User = require("../models/user")

const sessionConfig = {
    secret: process.env.SECRET,
    store: MongoStore.create({mongoUrl: process.env.DATABASE_URL}),
    resave: false,
    saveUninitialized: true
}
///////////////////////
// Custom Middelware
///////////////////////
const models = (req, res, next) => {
    req.models = {
        Entry,
        User
    }
    next()
}
///////////////////////////
// Export
///////////////////////////
module.exports = function(app) {
    app.use(express.urlencoded({extended: true}))
    app.use(express.json())
    app.use(methodOverride('_method'))
    app.use("/static", express.static('public'))
    app.use(morgan('tiny'))
    app.use(session(sessionConfig))
    app.use(models)
}