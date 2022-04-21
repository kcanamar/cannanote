////////////////////////
// Setup - Import deps and create app object
////////////////////////
require('dotenv').config()
require('./config/db')
const express = require('express')
const middleware = require('./middleware/mid')
const cannaRouter = require('./routes/entires')
const app = express()
const PORT = process.env.PORT
//////////////////////
// Declare Middleware
//////////////////////
middleware(app)
///////////////////////
// Declare Routes and Routers 
///////////////////////
app.use("/", cannaRouter)
///////////////////////////
// Server Listener
///////////////////////////
app.listen(PORT, () => console.log(`Who..Are..You?...${PORT}`))