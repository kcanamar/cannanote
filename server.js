////////////////////////
// Setup - Import deps and create app object
////////////////////////
require('dotenv').config()
const express = require('express')
const methodOverride = require('method-override')
const mongoose = require('mongoose')
const morgan = require('morgan')
const cannaRouter = require('./routes/entires')
const app = express()
const PORT = process.env.PORT || 3001
//////////////////////
// Declare Middleware
//////////////////////
app.use(methodOverride('_method'))
app.use("/static", express.static('public'))
app.use(express.urlencoded({extended: true}))
app.use(morgan('tiny'))
///////////////////////
// Declare Routes and Routers 
///////////////////////
// INDUCES - Index, New, Delete, Update, Create, Edit, Show
// Index
app.get("/", (req, res) => {
    res.send('It works')
})
// New
app.get("/new", (req, res) => {
    res.send('Something New')
})
// Edit
app.get("/:id/edit", (req, res) => {
    res.send(`here is what you asked to edit ${req.params.id}`)
})
// Show
app.get("/:id", (req, res) => {
    res.send(`Order for ${req.params.id}`)
})
// Update
app.put("/:id", (req, res) => {
    console.log(`here is what you wanted ${req.params.id}`)
    res.send(`It is Changed ${req.body}`)
})
// Create 
app.post("/", (req, res) => {
    res.send(`here is what you made ${req.body}`)
})
// Delete
app.delete("/:id", (req, res) => {
    res.send(`${req.params.id}....Target neutralized`)
})
///////////////////////////
// Server Listener
///////////////////////////
// MongoDB & Mongoose 
mongoose.connect(process.env.DATABASE_URL, {
    useNewUrlParser: true,
    useUnifiedTopology: true
})
const db = mongoose.connection;
db.on('error', (err) => console.log(err.message + "yeah... that didn't work"))
db.on('connected', () => console.log('mongoose connected'))
db.on('disconnected', () => console.log('mongoose disconnected'))

app.listen(PORT, () => console.log(`Who..Are..You?...${PORT}`))